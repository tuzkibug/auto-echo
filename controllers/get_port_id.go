package controllers

import (
	//"encoding/json"
	"fmt"
	"net/http"

	"bytes"
	"io/ioutil"

	"github.com/labstack/echo"
	//"github.com/tuzkibug/auto-echo/base"
)

//获取虚拟机端口ID

func Getportid(c echo.Context) (err error) {
	m := new(MsgPortMac)
	if err = c.Bind(m); err != nil {
		return
	}

	mac := m.PortMac

	url := "http://10.10.108.250:9696/v2.0/ports?mac_address=" + mac + "&fields=id"

	var jsonStr = []byte("")
	//fmt.Println("jsonStr", jsonStr)
	//fmt.Println("new_str", bytes.NewBuffer(jsonStr))

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	// req.Header.Set("X-Custom-Header", "myvalue")
	token := "gAAAAABeqUjaFEi2qKFJUTBRJR26LNEaVO-KZ7o0PCjk0YFUZZZNLuvEyC-3jQWgcyB8mBBdvRdK4qELv3z8Ew3NPABoUWlmetiC-VYbHAgU8Z8MLIcGzzh6sBNSXXoW41RPVqUokwkK0-eZRsFnBhgHrTbrNoqbFvNV2PsoTRFS7qnfqoZuSSk"
	req.Header.Set("X-Auth-Token", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	//fmt.Println("status", resp.Status)
	//fmt.Println("response:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body:", string(body))

	return
}
