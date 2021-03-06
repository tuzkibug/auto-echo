package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"

	//"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	log "github.com/sirupsen/logrus"
	"github.com/tuzkibug/auto-echo/base"
)

//自动化部署mysql主备集群

func BuilMysqlCluster(c echo.Context) (err error) {
	m := new(MsgMysqlCluster)
	if err = c.Bind(m); err != nil {
		log.Error(err)
		return
	}
	//openstack用户认证
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://" + m.OpenstackIP + ":5000/v3",
		Username:         m.Username,
		Password:         m.Password,
		DomainName:       m.DomainName,
		TenantID:         m.TenantID,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Error(err)
		return
	}

	//修改并生成master启动脚本
	base.ModifyMasterScript(m.VMRootPassword, m.MysqlRootPassword)
	//拉起主mysql虚拟机
	master_id := base.CreateMysql(provider, "master.txt", m.FlavorID, m.ImageID, m.NetworkID)
	//获取虚拟机IP,MAC
	i := 0
LOOP1:
	if i == 49 {
		log.Error("无法获取虚拟机信息，请检查虚拟机是否正常启动")
		return c.String(http.StatusNotFound, "无法获取虚拟机信息，请检查虚拟机是否正常启动")
	}
	i++
	master_ip := base.GetServerIP(provider, master_id)
	master_detail := *master_ip
	if master_detail.Status != "ACTIVE" {
		time.Sleep(5 * time.Second)
		goto LOOP1
	}
	master_addr := master_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	master_mac_addr := master_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]
	log.Info("Master ip is " + master_addr.(string))
	log.Info("Master mac is " + master_mac_addr.(string))

	//获取用户token
	username := m.Username
	password := m.Password
	domainname := m.DomainName
	url := "http://" + m.OpenstackIP + ":5000/v3/auth/tokens"
	reqbody := `{"auth": {"identity": {"methods": ["password"],"password": {"user": {"name": "` + username + `","domain": {"name": "` + domainname + `"},"password": "` + password + `"}}}}}`

	var jsonStr1 = []byte(reqbody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr1))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	token := resp.Header.Get("X-Subject-Token")

	//通过mac地址获取主mysql虚拟机port_id

	port_url := "http://" + m.OpenstackIP + ":9696/v2.0/ports?mac_address=" + master_mac_addr.(string) + "&fields=id"

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

	//绑定浮动IP
	//api地址/v2.0/floatingips
	//http://10.10.108.250:9696/v2.0/floatingips
	floating_url := "http://" + m.OpenstackIP + ":9696/v2.0/floatingips"
	floating_req_body := `{"floatingip": {"floating_network_id": "` + m.FloatingNetworkID + `","tenant_id": "` + m.TenantID + `","project_id": "` + m.TenantID + `","port_id": "` + port_id + `","fixed_ip_address": "` + master_addr.(string) + `"}}`

	var jsonStr3 = []byte(floating_req_body)
	req3, err := http.NewRequest("POST", floating_url, bytes.NewBuffer(jsonStr3))
	req3.Header.Set("X-Auth-Token", token)

	client3 := &http.Client{}
	resp3, err := client3.Do(req3)
	if err != nil {
		log.Error(err)
	}
	defer resp3.Body.Close()

	body3, _ := ioutil.ReadAll(resp3.Body)

	__fResponse := FIP{}
	if err := json.Unmarshal(body3, &__fResponse); err != nil {
		return err
	}

	//等待一段时间后，尝试连接数据库来确认是否安装完毕
	time.Sleep(100 * time.Second)
	db, err := sql.Open("mysql", "root:"+m.MysqlRootPassword+"@tcp("+__fResponse.FloatingIp.FloatingIp+":3306)/mysql?charset=utf8")
	if err != nil {
		log.Error("创建数据库对象失败")
		return
	}
	defer db.Close() // 延迟关闭 db对象创建成功后才可以调用close方法

	// 实际去尝试连接数据库
	for i := 0; i < 50; i++ {
		if i == 49 {
			log.Error("主节点连接异常，请检查")
			return c.String(http.StatusOK, "主节点连接异常，请检查")
		}
		err = nil
		err = db.Ping()
		if err == nil {
			log.Info("连接数据库主节点成功")
			break
		}
		log.Warn("暂无法连接数据库主节点，请稍后")
		time.Sleep(10 * time.Second)
	}

	//修改并生成slave启动脚本
	base.ModifySlaveScript(m.VMRootPassword, m.MysqlRootPassword, master_addr.(string))
	//拉起备mysql虚拟机
	slave_id := base.CreateMysql(provider, "slave.txt", m.FlavorID, m.ImageID, m.NetworkID)
	//获取虚拟机IP,MAC
	j := 0
LOOP2:
	if j == 49 {
		log.Error("无法获取虚拟机信息，请检查虚拟机是否正常启动")
		return c.String(http.StatusNotFound, "无法获取虚拟机信息，请检查虚拟机是否正常启动")
	}
	j++
	slave_ip := base.GetServerIP(provider, slave_id)
	slave_detail := *slave_ip
	if slave_detail.Status != "ACTIVE" {
		time.Sleep(5 * time.Second)
		goto LOOP2
	}
	slave_addr := slave_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	slave_mac_addr := slave_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]
	log.Info("Slave ip is " + slave_addr.(string))
	log.Info("Slave mac is " + slave_mac_addr.(string))
	log.Info("Mysql service is up in " + __fResponse.FloatingIp.FloatingIp + ":3306" + ". Username is root. Password is " + m.MysqlRootPassword)

	return c.String(http.StatusOK, __fResponse.FloatingIp.FloatingIp+" "+"3306"+" "+"root "+m.MysqlRootPassword)
}

func init() {
	// 设置日志格式为json格式
	log.SetFormatter(&log.TextFormatter{})

	// 设置将日志输出到标准输出（默认的输出为stderr，标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(os.Stdout)

	// 设置日志级别为warn以上
	log.SetLevel(log.InfoLevel)
}
