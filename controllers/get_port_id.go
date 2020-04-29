package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"bytes"
	"io/ioutil"

	"github.com/labstack/echo"
	//"github.com/tuzkibug/auto-echo/base"
)

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
	token := "gAAAAABep9PEo5wKYIAyD7l4rjoNGAZOKaGqWPBkmn4c0wCHZuz8fRTTEpfTo-8qqkWaFL9mmy2KMtR384-y5b0UhqIkgZAkJUQAiJD9eFY6WCUNQgInhYaKzdPb8lHmr3PrNiVWy1f2v8-B1dJ5ZIpN6a2ytGxAXl1AVVeSkk7JpqTOzcH2Ko4"
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
	//fmt.Println("response Body:", string(body))

	var dat []map[string]interface{}
	var ports map[string]interface{}

	ports["ports"] = dat
	if err := json.Unmarshal(body, &ports); err == nil {
		fmt.Println(dat)
	} else {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, dat)
}
