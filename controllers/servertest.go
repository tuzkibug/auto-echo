package controllers

import (
	"net/http"

	"fmt"
	"os"
	"strings"

	"github.com/labstack/echo"
	"github.com/tuzkibug/auto-echo/base"
)

func Servertest(c echo.Context) (err error) {
	filepath := "base_slave.txt"
	file, err := base.ReadAll(filepath)
	filestr := string(file)
	filestr_1 := strings.Replace(filestr, "aaa", "root", -1)
	filestr_2 := strings.Replace(filestr_1, "ccc", "root", -1)
	filestr_3 := strings.Replace(filestr_2, "ddd", "192.168.100.100", -1)

	fileName := "slave.txt"
	dstFile, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer dstFile.Close()
	s := filestr_3
	dstFile.WriteString(s)
	return c.String(http.StatusOK, filestr_3)
}
