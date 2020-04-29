package base

import (
	"fmt"
	"os"
	"strings"
)

func ModifyMasterScript(vmpassword string, mysqlpassword string) {
	filepath := "base.txt"
	file, err := ReadAll(filepath)
	filestr := string(file)
	filestr_1 := strings.Replace(filestr, "aaa", vmpassword, -1)
	filestr_2 := strings.Replace(filestr_1, "ccc", mysqlpassword, -1)

	fileName := "master.txt"
	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	s := filestr_2
	dstFile.WriteString(s)
	return
}

func ModifySlaveScript(vmpassword string, mysqlpassword string, masterip string) {
	filepath := "base_slave.txt"
	file, err := ReadAll(filepath)
	filestr := string(file)
	filestr_1 := strings.Replace(filestr, "aaa", vmpassword, -1)
	filestr_2 := strings.Replace(filestr_1, "ccc", mysqlpassword, -1)
	filestr_3 := strings.Replace(filestr_2, "ddd", masterip, -1)

	fileName := "slave.txt"
	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	s := filestr_3
	dstFile.WriteString(s)
	return
}
