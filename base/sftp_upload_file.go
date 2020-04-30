package base

import (
	"fmt"
	"io/ioutil"

	"os"
	"path"

	"github.com/pkg/sftp"
)

//SSH上传文件到指定路径

func UploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		return

	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)

	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		return

	}
	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		return

	}
	dstFile.Write(ff)
	fmt.Println(localFilePath + "  copy file to remote server finished!")
}
