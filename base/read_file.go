package base

import (
	"io/ioutil"
	"os"
)

//读文件到字符数组

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}
