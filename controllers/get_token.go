package controllers

import (
	"fmt"
	"net/http"

	"bytes"
	//"io/ioutil"

	"github.com/labstack/echo"
	//"github.com/tuzkibug/auto-echo/base"
)

func Getusertoken(c echo.Context) (err error) {
	m := new(MsgMysqlCreate)
	if err = c.Bind(m); err != nil {
		return
	}

	username := m.Username
	password := m.Password
	domainname := m.DomainName
	url := "http://10.10.108.250:5000/v3/auth/tokens"
	reqbody := "{\"auth\": {\"identity\": {\"methods\": [\"password\"],\"password\": {\"user\": {\"name\": \"" + username + "\",\"domain\": {\"name\": \"" + domainname + "\"},\"password\": \"" + password + "\"}}}}}"

	var jsonStr = []byte(reqbody)
	fmt.Println("jsonStr", jsonStr)
	fmt.Println("new_str", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	//fmt.Println("status", resp.Status)
	//fmt.Println("response:", resp.Header)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))

	return c.String(http.StatusOK, resp.Header.Get("X-Subject-Token"))
}
