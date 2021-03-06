package controllers

import (
	"bytes"
	"encoding/json"

	//"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/sftp"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	log "github.com/sirupsen/logrus"
	"github.com/tuzkibug/auto-echo/base"
)

//自动化部署CDH集群

func BuilCDHCluster(c echo.Context) (err error) {
	m := new(MsgCDHCluster)
	if err = c.Bind(m); err != nil {
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

	//拉起server虚拟机
	server_name := base.CreateCDHServerName() + "001"
	server_id := base.CreateCDHServer(provider, server_name, m.ServerFlavorID, m.ServerImageID, m.NetworkID)
	//拉起agent虚拟机
	a1_name := base.CreateCDHAgentName() + "001"
	agent1_id := base.CreateCDHAgent(provider, a1_name, m.AgentFlavorID, m.AgentImageID, m.NetworkID)
	a2_name := base.CreateCDHAgentName() + "002"
	agent2_id := base.CreateCDHAgent(provider, a2_name, m.AgentFlavorID, m.AgentImageID, m.NetworkID)
	a3_name := base.CreateCDHAgentName() + "003"
	agent3_id := base.CreateCDHAgent(provider, a3_name, m.AgentFlavorID, m.AgentImageID, m.NetworkID)

	//获取server虚拟机IP,MAC
	s_count := 0
LOOP0:
	if s_count == 49 {
		log.Error("无法获取虚拟机信息，请检查虚拟机是否正常启动")
		return c.String(http.StatusNotFound, "无法获取虚拟机信息，请检查虚拟机是否正常启动")
	}
	s_count++
	server_ip := base.GetServerIP(provider, server_id)
	server_detail := *server_ip
	if server_detail.Status != "ACTIVE" {
		log.Warn("等待虚拟机启动，请稍后")
		time.Sleep(5 * time.Second)
		goto LOOP0
	}

	server_addr := server_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	server_mac_addr := server_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]
	log.Info("Server ip is " + server_addr.(string))
	log.Info("Server mac is " + server_mac_addr.(string))

	//获取agent虚拟机IP,MAC
	a1_count := 0
LOOP1:
	if a1_count == 49 {
		log.Error("无法获取虚拟机信息，请检查虚拟机是否正常启动")
		return c.String(http.StatusNotFound, "无法获取虚拟机信息，请检查虚拟机是否正常启动")
	}
	a1_count++
	agent1_ip := base.GetServerIP(provider, agent1_id)
	agent1_detail := *agent1_ip
	if agent1_detail.Status != "ACTIVE" {
		log.Warn("等待虚拟机启动，请稍后")
		time.Sleep(5 * time.Second)
		goto LOOP1
	}
	agent1_addr := agent1_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	agent1_mac_addr := agent1_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]
	log.Info("Agent1 ip is " + agent1_addr.(string))
	log.Info("Agent1 mac is " + agent1_mac_addr.(string))

	a2_count := 0
LOOP2:
	if a2_count == 49 {
		log.Error("无法获取虚拟机信息，请检查虚拟机是否正常启动")
		return c.String(http.StatusNotFound, "无法获取虚拟机信息，请检查虚拟机是否正常启动")
	}
	a2_count++
	agent2_ip := base.GetServerIP(provider, agent2_id)
	agent2_detail := *agent2_ip
	if agent2_detail.Status != "ACTIVE" {
		log.Warn("等待虚拟机启动，请稍后")
		time.Sleep(5 * time.Second)
		goto LOOP2
	}
	agent2_addr := agent2_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	agent2_mac_addr := agent2_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]
	log.Info("Agent2 ip is " + agent2_addr.(string))
	log.Info("Agent2 mac is " + agent2_mac_addr.(string))

	a3_count := 0
LOOP3:
	if a3_count == 49 {
		log.Error("无法获取虚拟机信息，请检查虚拟机是否正常启动")
		return c.String(http.StatusNotFound, "无法获取虚拟机信息，请检查虚拟机是否正常启动")
	}
	a3_count++
	agent3_ip := base.GetServerIP(provider, agent3_id)
	agent3_detail := *agent3_ip
	if agent3_detail.Status != "ACTIVE" {
		log.Warn("等待虚拟机启动，请稍后")
		time.Sleep(5 * time.Second)
		goto LOOP3
	}
	agent3_addr := agent3_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	agent3_mac_addr := agent3_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]
	log.Info("Agent3 ip is " + agent3_addr.(string))
	log.Info("Agent3 mac is " + agent3_mac_addr.(string))

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

	//通过mac地址获取server虚拟机port_id

	port_url := "http://" + m.OpenstackIP + ":9696/v2.0/ports?mac_address=" + server_mac_addr.(string) + "&fields=id"

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
	log.Info("Server Port ID is " + port_id)

	//server绑定浮动IP
	//api地址/v2.0/floatingips
	//http://10.10.191.250:9696/v2.0/floatingips
	floating_url := "http://" + m.OpenstackIP + ":9696/v2.0/floatingips"
	floating_ip_network_id := m.FloatingNetworkID
	floating_req_body := `{"floatingip": {"floating_network_id": "` + floating_ip_network_id + `","tenant_id": "` + m.TenantID + `","project_id": "` + m.TenantID + `","port_id": "` + port_id + `","fixed_ip_address": "` + server_addr.(string) + `"}}`

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

	__serverfResponse := FIP{}
	if err := json.Unmarshal(body3, &__serverfResponse); err != nil {
		log.Error(err)
		return err
	}
	log.Info("Server Floating IP is " + __serverfResponse.FloatingIp.FloatingIp)

	//通过mac地址获取agent1虚拟机port_id
	port_url = "http://" + m.OpenstackIP + ":9696/v2.0/ports?mac_address=" + agent1_mac_addr.(string) + "&fields=id"
	jsonStr2 = []byte("")
	req2, err = http.NewRequest("GET", port_url, bytes.NewBuffer(jsonStr2))
	req2.Header.Set("X-Auth-Token", token)
	client2 = &http.Client{}
	resp2, err = client2.Do(req2)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
	body, _ = ioutil.ReadAll(resp2.Body)
	str = string(body)
	port_id_1 := str[17:53]
	log.Info("Agent1 Port ID is " + port_id_1)

	//agent1绑定浮动IP
	floating_req_body = `{"floatingip": {"floating_network_id": "` + floating_ip_network_id + `","tenant_id": "` + m.TenantID + `","project_id": "` + m.TenantID + `","port_id": "` + port_id_1 + `","fixed_ip_address": "` + agent1_addr.(string) + `"}}`
	jsonStr3 = []byte(floating_req_body)
	req3, err = http.NewRequest("POST", floating_url, bytes.NewBuffer(jsonStr3))
	req3.Header.Set("X-Auth-Token", token)
	client3 = &http.Client{}
	resp3, err = client3.Do(req3)
	if err != nil {
		panic(err)
	}
	defer resp3.Body.Close()
	body3, _ = ioutil.ReadAll(resp3.Body)
	__a1fResponse := FIP{}
	if err = json.Unmarshal(body3, &__a1fResponse); err != nil {
		log.Error(err)
		return err
	}
	log.Info("Agent1 Floating IP is " + __a1fResponse.FloatingIp.FloatingIp)

	//通过mac地址获取agent2虚拟机port_id
	port_url = "http://" + m.OpenstackIP + ":9696/v2.0/ports?mac_address=" + agent2_mac_addr.(string) + "&fields=id"
	jsonStr2 = []byte("")
	req2, err = http.NewRequest("GET", port_url, bytes.NewBuffer(jsonStr2))
	req2.Header.Set("X-Auth-Token", token)
	client2 = &http.Client{}
	resp2, err = client2.Do(req2)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
	body, _ = ioutil.ReadAll(resp2.Body)
	str = string(body)
	port_id_2 := str[17:53]
	log.Info("Agent2 Port ID is " + port_id_2)

	//agent2绑定浮动IP
	floating_req_body = `{"floatingip": {"floating_network_id": "` + floating_ip_network_id + `","tenant_id": "` + m.TenantID + `","project_id": "` + m.TenantID + `","port_id": "` + port_id_2 + `","fixed_ip_address": "` + agent2_addr.(string) + `"}}`
	jsonStr3 = []byte(floating_req_body)
	req3, err = http.NewRequest("POST", floating_url, bytes.NewBuffer(jsonStr3))
	req3.Header.Set("X-Auth-Token", token)
	client3 = &http.Client{}
	resp3, err = client3.Do(req3)
	if err != nil {
		panic(err)
	}
	defer resp3.Body.Close()
	body3, _ = ioutil.ReadAll(resp3.Body)
	__a2fResponse := FIP{}
	if err = json.Unmarshal(body3, &__a2fResponse); err != nil {
		log.Error(err)
		return err
	}
	log.Info("Agent2 Floating IP is " + __a2fResponse.FloatingIp.FloatingIp)

	//通过mac地址获取agent3虚拟机port_id
	port_url = "http://" + m.OpenstackIP + ":9696/v2.0/ports?mac_address=" + agent3_mac_addr.(string) + "&fields=id"
	jsonStr2 = []byte("")
	req2, err = http.NewRequest("GET", port_url, bytes.NewBuffer(jsonStr2))
	req2.Header.Set("X-Auth-Token", token)
	client2 = &http.Client{}
	resp2, err = client2.Do(req2)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
	body, _ = ioutil.ReadAll(resp2.Body)
	str = string(body)
	port_id_3 := str[17:53]
	log.Info("Agent3 Port ID is " + port_id_3)

	//agent3绑定浮动IP
	floating_req_body = `{"floatingip": {"floating_network_id": "` + floating_ip_network_id + `","tenant_id": "` + m.TenantID + `","project_id": "` + m.TenantID + `","port_id": "` + port_id_3 + `","fixed_ip_address": "` + agent3_addr.(string) + `"}}`
	jsonStr3 = []byte(floating_req_body)
	req3, err = http.NewRequest("POST", floating_url, bytes.NewBuffer(jsonStr3))
	req3.Header.Set("X-Auth-Token", token)
	client3 = &http.Client{}
	resp3, err = client3.Do(req3)
	if err != nil {
		panic(err)
	}
	defer resp3.Body.Close()
	body3, _ = ioutil.ReadAll(resp3.Body)
	__a3fResponse := FIP{}
	if err = json.Unmarshal(body3, &__a3fResponse); err != nil {
		log.Error(err)
		return err
	}
	log.Info("Agent3 Floating IP is " + __a3fResponse.FloatingIp.FloatingIp)

	//分别执行脚本
	cdhuser := "root"
	cdhpassword := "Admin123456"
	ciphers := []string{}

	//server删除hosts文件
	ss_count := 0
LOOP4:
	if ss_count == 49 {
		log.Error("无法连接至server虚拟机，请检查")
		return c.String(http.StatusNotFound, "无法连接至server虚拟机，请检查")
	}
	ss_count++
	session, err := base.Sshconnect(cdhuser, cdhpassword, __serverfResponse.FloatingIp.FloatingIp, "", 22, ciphers)
	if err != nil {
		log.Error(err)
		time.Sleep(5 * time.Second)
		goto LOOP4
	}
	defer session.Close()
	var serverstdoutBuf bytes.Buffer
	session.Stdout = &serverstdoutBuf
	session.Run("rm -rf /etc/hosts")
	log.Info(server_name + "删除初始/etc/hosts文件成功")

	//agent1删除hosts文件
	a1s_count := 0
LOOP5:
	if a1s_count == 49 {
		log.Error("无法连接至agent001虚拟机，请检查")
		return c.String(http.StatusNotFound, "无法连接至agent001虚拟机，请检查")
	}
	a1s_count++
	a1session, err := base.Sshconnect(cdhuser, cdhpassword, __a1fResponse.FloatingIp.FloatingIp, "", 22, ciphers)
	if err != nil {
		log.Error(err)
		time.Sleep(5 * time.Second)
		goto LOOP5
	}
	defer a1session.Close()
	var a1stdoutBuf bytes.Buffer
	a1session.Stdout = &a1stdoutBuf
	a1session.Run("rm -rf /etc/hosts")
	log.Info(a1_name + "删除初始/etc/hosts文件成功")

	//agent2删除hosts文件
	a2s_count := 0
LOOP6:
	if a2s_count == 49 {
		log.Error("无法连接至agent002虚拟机，请检查")
		return c.String(http.StatusNotFound, "无法连接至agent002虚拟机，请检查")
	}
	a2s_count++
	a2session, err := base.Sshconnect(cdhuser, cdhpassword, __a2fResponse.FloatingIp.FloatingIp, "", 22, ciphers)
	if err != nil {
		log.Error(err)
		time.Sleep(5 * time.Second)
		goto LOOP6
	}
	defer a2session.Close()
	var a2stdoutBuf bytes.Buffer
	a2session.Stdout = &a2stdoutBuf
	a2session.Run("rm -rf /etc/hosts")
	log.Info(a2_name + "删除初始/etc/hosts文件成功")

	//agent3删除hosts文件
	a3s_count := 0
LOOP7:
	if a3s_count == 49 {
		log.Error("无法连接至agent003虚拟机，请检查")
		return c.String(http.StatusNotFound, "无法连接至agent003虚拟机，请检查")
	}
	a3s_count++
	a3session, err := base.Sshconnect(cdhuser, cdhpassword, __a3fResponse.FloatingIp.FloatingIp, "", 22, ciphers)
	if err != nil {
		log.Error(err)
		time.Sleep(5 * time.Second)
		goto LOOP7
	}
	defer a3session.Close()
	var a3stdoutBuf bytes.Buffer
	a3session.Stdout = &a3stdoutBuf
	a3session.Run("rm -rf /etc/hosts")
	log.Info(a3_name + "删除初始/etc/hosts文件成功")

	//编辑etc/hosts文件
	base.ModifyEtcHosts(server_addr.(string), server_name, agent1_addr.(string), a1_name, agent2_addr.(string), a2_name, agent3_addr.(string), a3_name)
	//sftp上传编辑后的文件到server虚拟机
	var sftpClient *sftp.Client
	sftpClient, err = Connect(cdhuser, cdhpassword, __serverfResponse.FloatingIp.FloatingIp, 22)
	if err != nil {
		log.Error(err)
		return
	}
	defer sftpClient.Close()

	_, errStat := sftpClient.Stat("/etc/")
	if errStat != nil {
		log.Error(errStat)
		return
	}
	base.UploadFile(sftpClient, "hosts", "/etc/")

	//sftp上传编辑后的文件到a1虚拟机
	sftpClient, err = Connect(cdhuser, cdhpassword, __a1fResponse.FloatingIp.FloatingIp, 22)
	if err != nil {
		log.Error(err)
		return
	}
	defer sftpClient.Close()

	_, errStat = sftpClient.Stat("/etc/")
	if errStat != nil {
		log.Error(errStat)
		return
	}
	base.UploadFile(sftpClient, "hosts", "/etc/")

	//sftp上传编辑后的文件到a2虚拟机
	sftpClient, err = Connect(cdhuser, cdhpassword, __a2fResponse.FloatingIp.FloatingIp, 22)
	if err != nil {
		log.Error(err)
		return
	}
	defer sftpClient.Close()

	_, errStat = sftpClient.Stat("/etc/")
	if errStat != nil {
		log.Error(errStat)
		return
	}
	base.UploadFile(sftpClient, "hosts", "/etc/")

	//sftp上传编辑后的文件到a3虚拟机
	sftpClient, err = Connect(cdhuser, cdhpassword, __a3fResponse.FloatingIp.FloatingIp, 22)
	if err != nil {
		log.Error(err)
		return
	}
	defer sftpClient.Close()

	_, errStat = sftpClient.Stat("/etc/")
	if errStat != nil {
		log.Error(errStat)
		return
	}
	base.UploadFile(sftpClient, "hosts", "/etc/")

	//server执行安装脚本
	scmdstr := "/root/Config_CM_Server_arg.sh 1 " + server_name
	ss_count = 0
LOOP8:
	if ss_count == 49 {
		log.Error("无法连接至server虚拟机，请检查")
		return c.String(http.StatusNotFound, "无法连接至server虚拟机，请检查")
	}
	ss_count++
	session, err = base.Sshconnect(cdhuser, cdhpassword, __serverfResponse.FloatingIp.FloatingIp, "", 22, ciphers)
	if err != nil {
		log.Warn("server连接失败，请稍后")
		time.Sleep(5 * time.Second)
		goto LOOP8
	}
	defer session.Close()
	var serverstdoutBuf2 bytes.Buffer
	session.Stdout = &serverstdoutBuf2
	log.Info("This cmd will be executed " + scmdstr)
	session.Run(scmdstr)
	log.Info("server执行安装完成")

	//a1执行安装脚本
	a1cmdstr := "/root/Config_CM_Agent_arg.sh 1 " + server_name + " " + a1_name
	a1s_count = 0
LOOP9:
	if a1s_count == 49 {
		log.Error("无法连接至agent001虚拟机，请检查")
		return c.String(http.StatusNotFound, "无法连接至agent001虚拟机，请检查")
	}
	a1s_count++
	a1session, err = base.Sshconnect(cdhuser, cdhpassword, __a1fResponse.FloatingIp.FloatingIp, "", 22, ciphers)
	if err != nil {
		log.Warn("agent1连接失败，请稍后")
		time.Sleep(5 * time.Second)
		goto LOOP9
	}
	defer a1session.Close()
	var a1stdoutBuf2 bytes.Buffer
	a1session.Stdout = &a1stdoutBuf2
	log.Info("This cmd will be executed " + a1cmdstr)
	a1session.Run(a1cmdstr)
	log.Info("agent1执行安装完成")

	//a2执行安装脚本
	a2cmdstr := "/root/Config_CM_Agent_arg.sh 1 " + server_name + " " + a2_name
	a2s_count = 0
LOOP10:
	if a2s_count == 49 {
		log.Error("无法连接至agent002虚拟机，请检查")
		return c.String(http.StatusNotFound, "无法连接至agent002虚拟机，请检查")
	}
	a2s_count++
	a2session, err = base.Sshconnect(cdhuser, cdhpassword, __a2fResponse.FloatingIp.FloatingIp, "", 22, ciphers)
	if err != nil {
		log.Warn("agent2连接失败，请稍后")
		time.Sleep(5 * time.Second)
		goto LOOP10
	}
	defer a1session.Close()
	var a2stdoutBuf2 bytes.Buffer
	a2session.Stdout = &a2stdoutBuf2
	log.Info("This cmd will be executed " + a2cmdstr)
	a2session.Run(a2cmdstr)
	log.Info("agent2执行安装完成")

	//a3执行安装脚本
	a3cmdstr := "/root/Config_CM_Agent_arg.sh 1 " + server_name + " " + a3_name
	a3s_count = 0
LOOP11:
	if a3s_count == 49 {
		log.Error("无法连接至agent003虚拟机，请检查")
		return c.String(http.StatusNotFound, "无法连接至agent003虚拟机，请检查")
	}
	a3s_count++
	a3session, err = base.Sshconnect(cdhuser, cdhpassword, __a3fResponse.FloatingIp.FloatingIp, "", 22, ciphers)
	if err != nil {
		log.Warn("agent3连接失败，请稍后")
		time.Sleep(5 * time.Second)
		goto LOOP11
	}
	defer a3session.Close()
	var a3stdoutBuf2 bytes.Buffer
	a3session.Stdout = &a3stdoutBuf2
	log.Info("This cmd will be executed " + a3cmdstr)
	a3session.Run(a3cmdstr)
	log.Info("agent3执行安装完成")

	//等待重启完成，检查页面访问情况
	resp_count := 0
LOOP12:
	if resp_count == 49 {
		return c.String(http.StatusNotFound, "无法访问CM页面，请检查")
	}
	resp_count++
	testresp, err := http.Get("http://" + __serverfResponse.FloatingIp.FloatingIp + ":7180/cmf")
	if err != nil {
		log.Warn(err)
		time.Sleep(10 * time.Second)
		goto LOOP12
	}
	defer testresp.Body.Close()
	testbody, err := ioutil.ReadAll(testresp.Body)
	log.Info(string(testbody))
	log.Info("服务启动成功")
	log.Info("Service is up in " + __serverfResponse.FloatingIp.FloatingIp + ":7180")

	return c.String(http.StatusOK, __serverfResponse.FloatingIp.FloatingIp+":7180")

	//return c.String(http.StatusOK, server_addr.(string)+" "+server_mac_addr.(string)+" "+agent1_addr.(string)+" "+agent1_mac_addr.(string)+" "+agent2_addr.(string)+" "+agent2_mac_addr.(string)+" "+agent3_addr.(string)+" "+agent3_mac_addr.(string))
}
