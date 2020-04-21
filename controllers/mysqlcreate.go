package controllers

import (
	"net/http"
	//"encoding/json"
	"fmt"

	//  "io/ioutil"

	//	"time"

	"github.com/labstack/echo"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/tuzkibug/testopenstack/base"
)

func Createmysql(c echo.Context) (err error) {
	m := new(MsgMysqlCreate)
	//调用echo.Context的Bind函数将请求参数和User对象进行绑定。
	if err = c.Bind(m); err != nil {
		return
	}
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://10.10.100.55:5000/v3",
		Username:         m.Username,
		Password:         m.Password,
		DomainName:       m.DomainName,
		TenantID:         m.TenantID,
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	mysqlname := m.MysqlName
	base.CreateMysqlInstance(provider, mysqlname)

	return c.String(http.StatusOK, "The mysql server is created! The IP:PORT is "+""+":3306")

}