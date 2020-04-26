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
	"github.com/tuzkibug/auto-echo/base"
)

func MysqlIP(c echo.Context) (err error) {
	m := new(MsgMysqlDetail)
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
	server_id := m.MysqlID
	//base.CreateMysqlInstance(provider, mysqlname)
	server_ip := base.GetServerIP(provider, server_id)
	detail := *server_ip
	d := detail.Addresses["vxlan_172_16_2"].([]interface{})[0].(map[string]interface{})["addr"]

	return c.String(http.StatusOK, d.(string)+":3306")

}
