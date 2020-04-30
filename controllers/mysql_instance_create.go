package controllers

import (
	"net/http"

	"fmt"

	"github.com/labstack/echo"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/tuzkibug/auto-echo/base"
)

//测试用，拉起单独的mysql虚拟机

func Createmysql(c echo.Context) (err error) {
	m := new(MsgMysqlCreate)
	//调用echo.Context的Bind函数将请求参数和User对象进行绑定。
	if err = c.Bind(m); err != nil {
		return
	}
	opts := gophercloud.AuthOptions{
		IdentityEndpoint: "http://10.10.108.250:5000/v3",
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

	//base.CreateMysqlInstance(provider, mysqlname)
	file_name := "base_master.txt"
	flavor_id := "80588d70-7ba5-4863-8f77-d11170b2a007"
	image_id := "26e3fbd2-8beb-40fd-aa0f-dc285a56dcde"
	network_id := "2a8e355c-254e-4538-ab08-61a99c1da548"

	server_id := base.CreateMysql(provider, file_name, flavor_id, image_id, network_id)

	return c.String(http.StatusOK, server_id)

}
