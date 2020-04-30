package controllers

import (
	"fmt"

	"bytes"
	//"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

//测试用

func Servertest(c echo.Context) (err error) {
	m := new(MsgMysqlCluster)
	if err = c.Bind(m); err != nil {
		return
	}

	username := m.Username
	password := m.Password
	domainname := m.DomainName
	url := "http://10.10.108.250:5000/v3/auth/tokens"
	reqbody := "{\"auth\": {\"identity\": {\"methods\": [\"password\"],\"password\": {\"user\": {\"name\": \"" + username + "\",\"domain\": {\"name\": \"" + domainname + "\"},\"password\": \"" + password + "\"}}}}}"

	var jsonStr1 = []byte(reqbody)
	fmt.Println("jsonStr", jsonStr1)
	fmt.Println("new_str", bytes.NewBuffer(jsonStr1))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr1))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//fmt.Printf(reqbody)
	token := resp.Header.Get("X-Subject-Token")
	fmt.Println(token)

	mac := "fa:16:3e:aa:a4:8a"

	port_url := "http://10.10.108.250:9696/v2.0/ports?mac_address=" + mac + "&fields=id"
	fmt.Println(port_url)

	var jsonStr2 = []byte("")
	//fmt.Println("jsonStr", jsonStr)
	//fmt.Println("new_str", bytes.NewBuffer(jsonStr))

	req2, err := http.NewRequest("GET", port_url, bytes.NewBuffer(jsonStr2))
	// req.Header.Set("X-Custom-Header", "myvalue")

	req2.Header.Set("X-Auth-Token", token)

	client2 := &http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil {
		panic(err)
	}
	defer resp2.Body.Close()
	//fmt.Println("status", resp.Status)
	//fmt.Println("response:", resp.Header)
	body, _ := ioutil.ReadAll(resp2.Body)
	//fmt.Println("response Body:", string(body))
	str := string(body)
	port_id := str[17:53]

	fmt.Println(port_id)
	return
}
