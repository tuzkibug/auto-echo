package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/tuzkibug/auto-echo/base"
)

//自动化部署mysql主备集群

func BuilMysqlCluster(c echo.Context) (err error) {
	m := new(MsgMysqlCluster)
	if err = c.Bind(m); err != nil {
		return
	}
	//openstack用户认证
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://10.10.108.250:5000/v3",
		Username:         m.Username,
		Password:         m.Password,
		DomainName:       m.DomainName,
		TenantID:         m.TenantID,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	//修改并生成master启动脚本
	base.ModifyMasterScript(m.VMRootPassword, m.MysqlRootPassword)
	//拉起主mysql虚拟机
	master_id := base.CreateMysql(provider, "master.txt", m.FlavorID, m.ImageID, m.NetworkID)
	//获取虚拟机IP,MAC
	time.Sleep(20 * time.Second)
	master_ip := base.GetServerIP(provider, master_id)
	master_detail := *master_ip
	master_addr := master_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	master_mac_addr := master_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]
	//等待主节点安装配置完成
	time.Sleep(120 * time.Second)

	//修改并生成slave启动脚本
	base.ModifySlaveScript(m.VMRootPassword, m.MysqlRootPassword, master_addr.(string))
	//拉起备mysql虚拟机
	slave_id := base.CreateMysql(provider, "slave.txt", m.FlavorID, m.ImageID, m.NetworkID)
	//获取虚拟机IP,MAC
	time.Sleep(20 * time.Second)
	slave_ip := base.GetServerIP(provider, slave_id)
	slave_detail := *slave_ip
	slave_addr := slave_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	slave_mac_addr := slave_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]

	//获取用户token
	username := m.Username
	password := m.Password
	domainname := m.DomainName
	url := "http://10.10.108.250:5000/v3/auth/tokens"
	reqbody := "{\"auth\": {\"identity\": {\"methods\": [\"password\"],\"password\": {\"user\": {\"name\": \"" + username + "\",\"domain\": {\"name\": \"" + domainname + "\"},\"password\": \"" + password + "\"}}}}}"

	var jsonStr1 = []byte(reqbody)
	fmt.Println("jsonStr", jsonStr1)
	fmt.Println("new_str", bytes.NewBuffer(jsonStr1))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr1))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	token := resp.Header.Get("X-Subject-Token")
	fmt.Println(token)

	//通过mac地址获取主mysql虚拟机port_id
	//mac := "fa:16:3e:aa:a4:8a"

	port_url := "http://10.10.108.250:9696/v2.0/ports?mac_address=" + master_mac_addr.(string) + "&fields=id"
	fmt.Println(port_url)

	var jsonStr2 = []byte("")

	req2, err := http.NewRequest("GET", port_url, bytes.NewBuffer(jsonStr2))

	req2.Header.Set("X-Auth-Token", token)

	client2 := &http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()

	body, _ := ioutil.ReadAll(resp2.Body)

	str := string(body)

	port_id := str[17:53]

	fmt.Println(port_id)

	//绑定浮动IP
	//api地址/v2.0/floatingips
	//http://10.10.108.250:9696/v2.0/floatingips
	floating_url := "http://10.10.108.250:9696/v2.0/floatingips"
	floating_ip_network_id := "b9f41ba5-c37b-43dd-ad8b-e90ffe871a08"
	floating_req_body := `{"floatingip": {"floating_network_id": "` + floating_ip_network_id + `","tenant_id": "` + m.TenantID + `","project_id": "` + m.TenantID + `","port_id": "` + port_id + `","fixed_ip_address": "` + master_addr.(string) + `"}}`

	var jsonStr3 = []byte(floating_req_body)
	req3, err := http.NewRequest("POST", floating_url, bytes.NewBuffer(jsonStr3))
	req3.Header.Set("X-Auth-Token", token)

	client3 := &http.Client{}
	resp3, err := client3.Do(req3)
	if err != nil {
		panic(err)
	}
	defer resp3.Body.Close()

	body3, _ := ioutil.ReadAll(resp3.Body)

	str3 := string(body3)

	fmt.Println(str3)

	return c.String(http.StatusOK, master_addr.(string)+"  "+master_mac_addr.(string)+"  "+slave_addr.(string)+"  "+slave_mac_addr.(string))
}
