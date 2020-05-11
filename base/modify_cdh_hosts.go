package base

import (
	"fmt"
	"os"
	"strings"
)

// /etc/hosts修改函数

func ModifyEtcHosts(server_ip string, server_name string, a1_ip string, a1_name string, a2_ip string, a2_name string, a3_ip string, a3_name string) {
	filepath := "hosts_base"
	file, err := ReadAll(filepath)
	filestr := string(file)
	filestr = strings.Replace(filestr, "server_ip", server_ip, -1)
	filestr = strings.Replace(filestr, "server_name", server_name, -1)
	filestr = strings.Replace(filestr, "a1_ip", a1_ip, -1)
	filestr = strings.Replace(filestr, "a1_name", a1_name, -1)
	filestr = strings.Replace(filestr, "a2_ip", a2_ip, -1)
	filestr = strings.Replace(filestr, "a2_name", a2_name, -1)
	filestr = strings.Replace(filestr, "a3_ip", a3_ip, -1)
	filestr = strings.Replace(filestr, "a3_name", a3_name, -1)

	fileName := "hosts"
	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	s := filestr
	dstFile.WriteString(s)
	return
}
