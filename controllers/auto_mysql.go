package controllers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/rackspace/gophercloud"
	"github.com/rackspace/gophercloud/openstack"
	"github.com/tuzkibug/auto-echo/base"
)

func BuilMysqlCluster(c echo.Context) (err error) {
	m := new(MsgMysqlCluster)
	if err = c.Bind(m); err != nil {
		return
	}
	//openstack用户认证
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

	//修改并生成master启动脚本
	base.ModifyMasterScript(m.VMRootPassword, m.MysqlRootPassword)
	//拉起主mysql虚拟机
	master_id := base.CreateMysql(provider, "master.txt", m.FlavorID, m.ImageID, m.NetworkID)
	//获取虚拟机IP,MAC
	master_ip := base.GetServerIP(provider, master_id)
	master_detail := *master_ip
	master_addr := master_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	master_mac_addr := master_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]

	//修改并生成slave启动脚本
	base.ModifySlaveScript(m.VMRootPassword, m.MysqlRootPassword, master_addr.(string))
	//拉起备mysql虚拟机
	slave_id := base.CreateMysql(provider, "slave.txt", m.FlavorID, m.ImageID, m.NetworkID)
	//获取虚拟机IP,MAC
	slave_ip := base.GetServerIP(provider, slave_id)
	slave_detail := *slave_ip
	slave_addr := slave_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["addr"]
	slave_mac_addr := slave_detail.Addresses[m.NetworkName].([]interface{})[0].(map[string]interface{})["OS-EXT-IPS-MAC:mac_addr"]

	return c.String(http.StatusOK, master_addr.(string)+"  "+master_mac_addr.(string)+"  "+slave_addr.(string)+"  "+slave_mac_addr.(string))
}
