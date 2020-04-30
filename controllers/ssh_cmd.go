package controllers

import (
	"bytes"
	"net/http"

	"github.com/labstack/echo"
	"github.com/tuzkibug/auto-echo/base"
)

//ssh连接虚拟机并执行shell命令或脚本

func SSH_run_cmd(c echo.Context) (err error) {
	u := new(MsgVMSSH)
	ciphers := []string{}
	if err = c.Bind(u); err != nil {
		return
	}
	session, err := base.Sshconnect(u.Username, u.Password, u.SshIP, "", u.Sshport, ciphers)
	if err != nil {
		return
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(u.Cmd)
	return c.String(http.StatusOK, stdoutBuf.String())
}
