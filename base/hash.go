package base

import (
	"crypto/md5"
	"fmt"
	"io"
	"time"
)

/*
随机生成哈希字符串
*/

func CreateRandom() string {
	t := time.Now()
	h := md5.New()
	io.WriteString(h, "PCL")
	io.WriteString(h, t.String())
	passwd := fmt.Sprintf("%x", h.Sum(nil))
	return passwd
}
