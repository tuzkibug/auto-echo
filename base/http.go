package base

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func httpDo(method string, url string, msg string) (result string) {
	client := &http.Client{}
	body := bytes.NewBuffer([]byte(msg))
	req, err := http.NewRequest(method,
		url,
		body)
	if err != nil {
		// handle error
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	result_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	result = string(result_body)
	return result
}

// post方式
func HttpDoPost(url string, msg string) (result string) {
	httpDo("POST", url, msg)
	return
}

// get方式
func HttpDoGet(url string, msg string) (result string) {
	httpDo("GET", url, msg)
	return
}
