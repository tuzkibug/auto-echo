package ob

import (
	"bytes"
	"encoding/json"
	"os"
	"strconv"

	//"fmt"
	"io/ioutil"
	"net/http"

	"time"

	"github.com/labstack/echo"
	//"github.com/pkg/sftp"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	log "github.com/sirupsen/logrus"
	"github.com/tuzkibug/auto-echo/base"
	"github.com/tuzkibug/auto-echo/controllers"
)

//定义CDH虚拟机对象，server和agent均属于该对象的实例
type CDHVM struct {
	role     string
	username string
	password string
	ip       string
	name     string
	fip      string
}

//对象方法：追加文件信息
func (cdh *CDHVM) AddInfo() {
	str := []byte("\n" + cdh.ip + " " + cdh.name)

	// 以追加模式打开文件
	txt, err := os.OpenFile(`hosts`, os.O_APPEND, 0666)

	defer txt.Close()
	if err != nil {
		panic(err)
	}

	// 写入文件
	n, err := txt.Write(str)
	// 当 n != len(b) 时，返回非零错误
	if err == nil && n != len(str) {
		log.Error(`错误代码：`, n)
		panic(err)
	}
}

//对象方法：上传hosts新文件
func (cdh *CDHVM) TransHosts() {
	ciphers := []string{}
	ss_count := 0
LOOP4:
	if ss_count == 49 {
		log.Error("无法连接至server虚拟机，请检查")
		return
	}
	ss_count++
	session, err := base.Sshconnect(cdh.username, cdh.password, cdh.fip, "", 22, ciphers)
	if err != nil {
		log.Error(err)
		time.Sleep(5 * time.Second)
		goto LOOP4
	}
	defer session.Close()
	var serverstdoutBuf bytes.Buffer
	session.Stdout = &serverstdoutBuf
	session.Run("rm -rf /etc/hosts")
	log.Info(cdh.name + "删除初始/etc/hosts文件成功")

	sftpClient, err := controllers.Connect(cdh.username, cdh.password, cdh.fip, 22)
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
}

//对象方法：创建server虚拟机
func (dd *CDHCluster) CreateServerVM(provider *gophercloud.ProviderClient, no int, id chan string) {
	server_name := base.CreateCDHServerName() + strconv.Itoa(no)
	server_id := base.CreateCDHServer(provider, server_name, dd.ServerFlavorID, dd.ServerImageID, dd.NetworkID)
	id <- server_id
}

//对象方法：创建agent虚拟机
func (dd *CDHCluster) CreateAgentVM(provider *gophercloud.ProviderClient, no int, id chan string) {
	a_name := base.CreateCDHAgentName() + strconv.Itoa(no)
	agent_id := base.CreateCDHAgent(provider, a_name, dd.AgentFlavorID, dd.AgentImageID, dd.NetworkID)
	id <- agent_id
}

//对象方法：配置浮动IP
func (dd *CDHCluster) SetFIP(server_ip string, server_mac string) (string, error) {
	username := dd.Username
	password := dd.Password
	domainname := dd.DomainName
	url := "http://" + dd.OpenstackIP + ":5000/v3/auth/tokens"
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
	port_url := "http://" + dd.OpenstackIP + ":9696/v2.0/ports?mac_address=" + server_mac + "&fields=id"
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

	floating_url := "http://" + dd.OpenstackIP + ":9696/v2.0/floatingips"
	floating_ip_network_id := dd.FloatingNetworkID
	floating_req_body := `{"floatingip": {"floating_network_id": "` + floating_ip_network_id + `","tenant_id": "` + dd.TenantID + `","project_id": "` + dd.TenantID + `","port_id": "` + port_id + `","fixed_ip_address": "` + server_ip + `"}}`

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
		return "no fip", err
	}
	return __serverfResponse.FloatingIp.FloatingIp, nil
}

func BuildCDHCluster(c echo.Context) (err error) {
	d := new(CDHCluster)
	if err = c.Bind(d); err != nil {
		return
	}
	//openstack用户认证
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://" + d.OpenstackIP + ":5000/v3",
		Username:         d.Username,
		Password:         d.Password,
		DomainName:       d.DomainName,
		TenantID:         d.TenantID,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		log.Error(err)
		return
	}

	//定义数组，存放server主机ID等信息
	var ids1 []string
	var names1 []string
	var ips1 []string
	var fips1 []string
	//并发创建server虚拟机并记录id
	idschs1 := make([]chan string, d.SeverNum)
	for i := 0; i < d.SeverNum; i++ {
		idschs1[i] = make(chan string)
		go d.CreateServerVM(provider, i, idschs1[i])
	}

	for _, idch1 := range idschs1 {
		ids1 = append(ids1, <-idch1)
	}

	i := 0
LOOP0:
	s_count := 0
LOOP1:
	if s_count == 49 {
		log.Error("无法获取虚拟机" + ids1[i] + "信息，请检查虚拟机是否正常启动")
		return c.String(http.StatusNotFound, "无法获取虚拟机"+ids1[i]+"信息，请检查虚拟机是否正常启动")
	}
	s_count++

	server := base.GetServerIP(provider, ids1[i])
	server_detail := *server
	if server_detail.Status != "ACTIVE" {
		log.Warn("等待虚拟机" + ids1[i] + "启动，请稍后")
		time.Sleep(5 * time.Second)
		goto LOOP1
	}
	server_name := server_detail.Name
	server_ip := server_detail.Addresses[d.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	server_mac := server_detail.Addresses[d.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]
	server_fip, _ := d.SetFIP(server_ip.(string), server_mac.(string))
	names1 = append(names1, server_name)
	ips1 = append(ips1, server_ip.(string))
	fips1 = append(fips1, server_fip)
	log.Info(server_name + " is active now. The ip is " + server_ip.(string) + ", floating ip is " + server_fip)

	i++
	if i < d.SeverNum {
		goto LOOP0
	}

	//定义数组，存放agent主机ID等信息
	var ids2 []string
	var names2 []string
	var ips2 []string
	var fips2 []string
	//并发创建agent虚拟机并记录id
	idschs2 := make([]chan string, d.AgentNum)
	for j := 0; j < d.AgentNum; j++ {
		idschs2[j] = make(chan string)
		go d.CreateAgentVM(provider, j, idschs2[j])
	}

	for _, idch2 := range idschs2 {
		ids2 = append(ids2, <-idch2)
	}

	j := 0
LOOP2:
	a_count := 0
LOOP3:
	if a_count == 49 {
		log.Error("无法获取虚拟机" + ids2[j] + "信息，请检查虚拟机是否正常启动")
		return c.String(http.StatusNotFound, "无法获取虚拟机"+ids2[j]+"信息，请检查虚拟机是否正常启动")
	}
	a_count++

	agent := base.GetServerIP(provider, ids2[j])
	agent_detail := *agent
	if agent_detail.Status != "ACTIVE" {
		log.Warn("等待虚拟机" + ids2[j] + "启动，请稍后")
		time.Sleep(5 * time.Second)
		goto LOOP3
	}
	agent_name := agent_detail.Name
	agent_ip := agent_detail.Addresses[d.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	agent_mac := agent_detail.Addresses[d.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]
	agent_fip, _ := d.SetFIP(agent_ip.(string), agent_mac.(string))
	names2 = append(names2, agent_name)
	ips2 = append(ips2, agent_ip.(string))
	fips2 = append(fips2, agent_fip)
	log.Info(agent_name + " is active now. The ip is " + agent_ip.(string) + ", floating ip is " + agent_fip)

	j++
	if j < d.AgentNum {
		goto LOOP2
	}

	//修改hosts文件，先拷贝，传完再删除
	input, err := ioutil.ReadFile("hosts_base")
	if err != nil {
		log.Error(err)
		return
	}

	err = ioutil.WriteFile("hosts", input, 0644)
	if err != nil {
		log.Error("Error creating", "hosts")
		log.Error(err)
		return
	}

	//追加server信息
	for i = 0; i < d.SeverNum; i++ {
		cdhserver := CDHVM{role: "server", username: "root", password: "Admin123456", ip: ips1[i], name: names1[i], fip: fips1[i]}
		cdhserver.AddInfo()
	}

	//追加agent信息
	for j = 0; j < d.AgentNum; j++ {
		cdhagent := CDHVM{role: "agent", username: "root", password: "Admin123456", ip: ips2[j], name: names2[j], fip: fips2[j]}
		cdhagent.AddInfo()
	}

	//server删除hosts文件并上传新文件
	for i = 0; i < d.SeverNum; i++ {
		cdhserver := CDHVM{role: "server", username: "root", password: "Admin123456", ip: ips1[i], name: names1[i], fip: fips1[i]}
		cdhserver.TransHosts()
	}

	//agent删除hosts文件并上传新文件
	for j = 0; j < d.AgentNum; j++ {
		cdhagent := CDHVM{role: "agent", username: "root", password: "Admin123456", ip: ips2[j], name: names2[j], fip: fips2[j]}
		cdhagent.TransHosts()
	}

	//删除本地hosts文件
	err = os.Remove("hosts")
	if err != nil {
		log.Error(err)
		return err
	}

	//所有虚拟机并行执行脚本

	return c.String(http.StatusOK, "OK")
}
