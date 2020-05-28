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

//对象方法:修改脚本
func (mm *MysqlCluster) ModifyScript(role string, vmrootpassword string, mysqlrootpassword string, masterip string) {
	if role == "master" {
		base.ModifyMasterScript(mm.VMRootPassword, mm.MysqlRootPassword)
	}
	if role == "slave" {
		base.ModifySlaveScript(mm.VMRootPassword, mm.MysqlRootPassword, masterip)
	}
}

//对象方法：创建虚机
func (mm *MysqlCluster) CreateVM(role string, provider *gophercloud.ProviderClient) string {
	if role == "master" {
		serverId := base.CreateMysql(provider, "master.txt", mm.FlavorID, mm.ImageID, mm.NetworkID)
		return serverId
	}
	if role == "slave" {
		slaveId := base.CreateMysql(provider, "slave.txt", mm.FlavorID, mm.ImageID, mm.NetworkID)
		return slaveId
	}
	return "create failed"
}

//对象方法：获取IP和MAC
func (mm *MysqlCluster) Getinfo(role string, t string, provider *gophercloud.ProviderClient, id string) string {
	i := 0
LOOP1:
	if i == 49 {
		log.Error("无法获取虚拟机信息，请检查虚拟机是否正常启动")
	}
	i++
	ip := base.GetServerIP(provider, id)
	detail := *ip
	if detail.Status != "ACTIVE" {
		time.Sleep(5 * time.Second)
		goto LOOP1
	}
	if role == "master" && t == "ip" {
		masterAddr := detail.Addresses[mm.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
		return masterAddr.(string)
	}
	if role == "master" && t == "mac" {
		masterMacAddr := detail.Addresses[mm.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]
		return masterMacAddr.(string)
	}
	if role == "slave" && t == "ip" {
		slaveAddr := detail.Addresses[mm.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
		return slaveAddr.(string)
	}
	if role == "slave" && t == "mac" {
		slaveMacAddr := detail.Addresses[mm.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]
		return slaveMacAddr.(string)
	}
	return "no info"
}

//对象方法：获取用户token
func (mm *MysqlCluster) GetToken(openstackIP string, username string, password string, domainname string) string {
	url := "http://" + openstackIP + ":5000/v3/auth/tokens"
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
	return resp.Header.Get("X-Subject-Token")
}

//对象方法：配置浮动IP
func (mm *MysqlCluster) SetFIP(token string, ip string, mac string) (string, error) {
	//获取port_id
	portUrl := "http://" + mm.OpenstackIP + ":9696/v2.0/ports?mac_address=" + mac + "&fields=id"

	var jsonStr = []byte("")

	req, err := http.NewRequest("GET", portUrl, bytes.NewBuffer(jsonStr))

	req.Header.Set("X-Auth-Token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	str := string(body)

	portId := str[17:53]

	//绑定浮动IP
	//api地址/v2.0/floatingips
	//http://10.10.108.250:9696/v2.0/floatingips
	floatingUrl := "http://" + mm.OpenstackIP + ":9696/v2.0/floatingips"
	floatingReqBody := `{"floatingip": {"floating_network_id": "` + mm.FloatingNetworkID + `","tenant_id": "` + mm.TenantID + `","project_id": "` + mm.TenantID + `","port_id": "` + portId + `","fixed_ip_address": "` + ip + `"}}`

	var jsonStr3 = []byte(floatingReqBody)
	req3, err := http.NewRequest("POST", floatingUrl, bytes.NewBuffer(jsonStr3))
	req3.Header.Set("X-Auth-Token", token)

	client3 := &http.Client{}
	resp3, err := client3.Do(req3)
	if err != nil {
		log.Error(err)
	}
	defer resp3.Body.Close()

	body3, _ := ioutil.ReadAll(resp3.Body)

	f := FIP{}
	if err := json.Unmarshal(body3, &f); err != nil {
		return "no FIP", err
	}

	return f.FloatingIp.FloatingIp, nil
}

//对象方法：连接数据库测试
func (mm *MysqlCluster) LinkTest(fip string) error {
	db, err := sql.Open("mysql", "root:"+mm.MysqlRootPassword+"@tcp("+fip+":3306)/mysql?charset=utf8")
	if err != nil {
		log.Error("创建数据库对象失败")
		return err
	}
	defer db.Close()

	for i := 0; i < 50; i++ {
		if i == 49 {
			log.Error("主节点连接异常，请检查")
			return err
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
	return nil

}

//自动化部署mysql主备集群
//主函数
func BuildMysqlCluster(c echo.Context) (err error) {
	m := new(MysqlCluster)
	if err = c.Bind(m); err != nil {
		log.Error(err)
		return
	}

	//openstack统一认证
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

	//修改master脚本
	m.ModifyScript("master", m.VMRootPassword, m.MysqlRootPassword, "")
	//拉起主mysql虚拟机并获取ID
	masterId := m.CreateVM("master", provider)
	//获取server虚拟机IP,MAC
	masterAddr := m.Getinfo("master", "ip", provider, masterId)
	masterMacAddr := m.Getinfo("master", "mac", provider, masterId)
	log.Info("Master ip is " + masterAddr)
	log.Info("Master mac is " + masterMacAddr)

	//获取用户token

	token := m.GetToken(m.OpenstackIP, m.Username, m.Password, m.DomainName)

	//配置浮动IP
	fip, _ := m.SetFIP(token, masterAddr, masterMacAddr)

	//等待一段时间后，尝试连接数据库来确认是否安装完毕
	time.Sleep(100 * time.Second)
	err = m.LinkTest(fip)
	if err != nil {
		return
	}

	//修改并生成slave启动脚本
	m.ModifyScript("slave", m.VMRootPassword, m.MysqlRootPassword, masterAddr)
	//拉起备mysql虚拟机
	slaveId := m.CreateVM("slave", provider)
	//获取slave虚拟机IP,MAC
	slaveAddr := m.Getinfo("slave", "ip", provider, slaveId)
	slaveMacAddr := m.Getinfo("slave", "mac", provider, slaveId)
	log.Info("Slave ip is " + slaveAddr)
	log.Info("Slave mac is " + slaveMacAddr)
	log.Info("Mysql service is up in " + fip + ":3306" + ". Username is root. Password is " + m.MysqlRootPassword)

	//删除本地master/slave临时文件
	err = os.Remove("master.txt")
	if err != nil {
		log.Error(err)
		return err
	}
	err = os.Remove("slave.txt")
	if err != nil {
		log.Error(err)
		return err
	}

	return c.String(http.StatusOK, fip+" "+"3306"+" "+"root "+m.MysqlRootPassword)
}
