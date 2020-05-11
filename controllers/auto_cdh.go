package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/tuzkibug/auto-echo/base"
)

//自动化部署mysql主备集群

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
		fmt.Printf("%v", err)
		return
	}

	//拉起server虚拟机
	server_id := base.CreateCDHServer(provider, m.FlavorID, m.ServerImageID, m.NetworkID)
	//拉起agent虚拟机
	agent1_id := base.CreateCDHAgent(provider, m.FlavorID, m.AgentImageID, m.NetworkID)
	agent2_id := base.CreateCDHAgent(provider, m.FlavorID, m.AgentImageID, m.NetworkID)
	agent3_id := base.CreateCDHAgent(provider, m.FlavorID, m.AgentImageID, m.NetworkID)

	//获取server虚拟机IP,MAC
LOOP0:
	server_ip := base.GetServerIP(provider, server_id)
	server_detail := *server_ip
	if server_detail.Status != "ACTIVE" {
		time.Sleep(5 * time.Second)
		goto LOOP0
	}

	server_addr := server_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	server_mac_addr := server_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]

	//获取agent虚拟机IP,MAC
LOOP1:
	agent1_ip := base.GetServerIP(provider, agent1_id)
	agent1_detail := *agent1_ip
	if agent1_detail.Status != "ACTIVE" {
		time.Sleep(5 * time.Second)
		goto LOOP1
	}
	agent1_addr := agent1_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	agent1_mac_addr := agent1_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]

LOOP2:
	agent2_ip := base.GetServerIP(provider, agent2_id)
	agent2_detail := *agent2_ip
	if agent2_detail.Status != "ACTIVE" {
		time.Sleep(5 * time.Second)
		goto LOOP2
	}
	agent2_addr := agent2_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	agent2_mac_addr := agent2_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]

LOOP3:
	agent3_ip := base.GetServerIP(provider, agent3_id)
	agent3_detail := *agent3_ip
	if agent3_detail.Status != "ACTIVE" {
		time.Sleep(5 * time.Second)
		goto LOOP3
	}
	agent3_addr := agent3_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	agent3_mac_addr := agent3_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]

	//获取用户token
	username := m.Username
	password := m.Password
	domainname := m.DomainName
	url := "http://" + m.OpenstackIP + ":5000/v3/auth/tokens"
	reqbody := "{\"auth\": {\"identity\": {\"methods\": [\"password\"],\"password\": {\"user\": {\"name\": \"" + username + "\",\"domain\": {\"name\": \"" + domainname + "\"},\"password\": \"" + password + "\"}}}}}"

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
	fmt.Println(str)

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

	str3 := string(body3)

	location := strings.IndexAny(str3, "floating_ip_address")

	__snsOauth2Response := FIP{}
	if err := json.Unmarshal(body3, &__snsOauth2Response); err != nil {
		return err
	}

	fmt.Println(__snsOauth2Response.FloatingIp.FloatingIp)

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
	str3 = string(body3)
	location = strings.IndexAny(str3, "floating_ip_address")
	fmt.Println(str3[location+23])

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
	str3 = string(body3)
	location = strings.IndexAny(str3, "floating_ip_address")
	fmt.Println(str3[location+23])

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
	str3 = string(body3)
	location = strings.IndexAny(str3, "floating_ip_address")
	fmt.Println(str3[location+23])

	return c.String(http.StatusOK, server_addr.(string)+" "+server_mac_addr.(string)+" "+agent1_addr.(string)+" "+agent1_mac_addr.(string)+" "+agent2_addr.(string)+" "+agent2_mac_addr.(string)+" "+agent3_addr.(string)+" "+agent3_mac_addr.(string))
}
