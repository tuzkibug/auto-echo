package controllers

import (
	"github.com/tuzkibug/auto-echo/base"

	"fmt"

	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

//ssh上传文件

func Connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		Timeout:         30 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //ssh.FixedHostKey(hostKey),
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func UploadSSH(c echo.Context) (err error) {
	u := new(MsgUploadSSH)
	if err = c.Bind(u); err != nil {
		return
	}

	var sftpClient *sftp.Client
	start := time.Now()
	sftpClient, err = Connect(u.Username, u.Password, u.SshIP, u.Sshport)
	if err != nil {
		return
	}
	defer sftpClient.Close()

	_, errStat := sftpClient.Stat(u.Remotepath)
	if errStat != nil {
		return
	}

	base.UploadFile(sftpClient, u.Localpath, u.Remotepath)
	elapsed := time.Since(start)
	fmt.Println("elapsed time : ", elapsed)
	return c.JSON(http.StatusOK, u)
}
