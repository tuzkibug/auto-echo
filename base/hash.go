package base

import (
	"crypto/md5"
	"fmt"
	"io"
	"math/rand"
	"time"
)

/*
随机生成哈希字符串作为mysql虚拟机名称
*/

func CreateRandom() string {
	t := time.Now()
	h := md5.New()
	io.WriteString(h, "PCL")
	io.WriteString(h, t.String())
	passwd := fmt.Sprintf("%x", h.Sum(nil))
	return passwd
}

/*
随机生成CDH要求的主机名称，仅小写字母和数字，带server和agent标识
*/
func CreateCDHServerName() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 7; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return "cdhserver" + string(result)
}

func CreateCDHAgentName() string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 7; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return "cdhagent" + string(result)
}
