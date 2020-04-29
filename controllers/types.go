package controllers

type MsgMysqlCreate struct {
	Username   string `json:"username" form:"username" query:"username"`
	Password   string `json:"password" form:"password" query:"password"`
	DomainName string `json:"domain_name" form:"domain_name" query:"domain_name"`
	TenantID   string `json:"tenant_id" form:"tenant_id" query:"tenant_id"`
	MysqlName  string `json:"mysql_name" form:"mysql_name" query:"mysql_name"`
}

/*
创建mysql虚拟机时构造的请求
header:
Content-Type:application/json
body:
{
  "username":"pcl",
  "password":"pcl@123",
  "domain_name":"default",
  "tenant_id":"6e57ee69fb0740fc89e53f3bea47a545",
  "mysql_name":"mymy"
}
*/

type MsgMysqlPasswordInitial struct {
	MysqlIP string `json:"mysql_ip" form:"mysql_ip" query:"mysql_ip"`
	//默认端口3306
	//MysqlPort   string `json:"mysql_port" form:"mysql_port" query:"mysql_port"`
	Newpassword string `json:"newpassword" form:"newpassword" query:"newpassword"`
}

/*
初始化mysql密码时构造的请求
{
  "mysql_ip":"192.168.56.109",
  "newpassword":"root"
}
*/

type MsgMysqlDetail struct {
	Username   string `json:"username" form:"username" query:"username"`
	Password   string `json:"password" form:"password" query:"password"`
	DomainName string `json:"domain_name" form:"domain_name" query:"domain_name"`
	TenantID   string `json:"tenant_id" form:"tenant_id" query:"tenant_id"`
	MysqlID    string `json:"mysql_id" form:"mysql_id" query:"mysql_id"`
}

/*
获取mysql IP的请求
header:
Content-Type:application/json
body:
{
  "username":"pcl",
  "password":"pcl@123",
  "domain_name":"default",
  "tenant_id":"6e57ee69fb0740fc89e53f3bea47a545",
  "mysql_id":"c6f645a5-5d5b-4e41-b89a-6b48dd1a23c5"
}
*/

type MsgVMSSH struct {
	Username string `json:"username" form:"username" query:"username"`
	Password string `json:"password" form:"password" query:"password"`
	SshIP    string `json:"sship" form:"sship" query:"sship"`
	Sshport  int    `json:"sshport" form:"sshport" query:"sshport"`
	Cmd      string `json:"cmd" form:"cmd" query:"cmd"`
}

/*
向远程主机发起SSH连接并执行命令时构造的请求
header:
Content-Type:application/json
body:
{
  "username":"root",
  "password":"root",
  "sship":"192.168.56.109",
  "sshport":22,
  "cmd":"/root/install_mysql.sh"
}
*/

type MsgUploadSSH struct {
	Username   string `json:"username" form:"username" query:"username"`
	Password   string `json:"password" form:"password" query:"password"`
	SshIP      string `json:"sship" form:"sship" query:"sship"`
	Sshport    int    `json:"sshport" form:"sshport" query:"sshport"`
	Localpath  string `json:"localpath" form:"localpath" query:"localpath"`
	Remotepath string `json:"remotepath" form:"remotepath" query:"remotepath"`
}

/*
向远程主机上传本地文件时构造的请求
header:
Content-Type:application/json
body:
{
  "username":"root",
  "password":"root",
  "sship":"192.168.56.109",
  "sshport":22,
  "localpath":"my.cnf.master",
  "remotepath":"/etc/"
}
*/

type MsgOPSDomain struct {
	DomainName string `json:"name" form:"name" query:"name"`
}

type MsgOPSUser struct {
	UserName     string `json:"name" form:"name" query:"name"`
	MsgOPSDomain `json:"domain" form:"domain" query:"domain"`
	UserPassword string `json:"password" form:"password" query:"password"`
}

type MsgOPSPassword struct {
	MsgOPSUser `json:"user" form:"user" query:"user"`
}

type MsgOPSIdentity struct {
	MsgOPSPassword `json:"identity" form:"identity" query:"identity"`
}

type MsgOPSAuth struct {
	MsgOPSIdentity `json:"auth" form:"auth" query:"auth"`
}

type MsgPortMac struct {
	PortMac string `json:"port_mac" form:"port_mac" query:"port_mac"`
}

/*
获取port id时构造的请求
body:
{
  "port_mac":"fa:16:3e:f4:48:f1"
}
*/

type MsgMysqlCluster struct {
	Username          string `json:"username" form:"username" query:"username"`
	Password          string `json:"password" form:"password" query:"password"`
	DomainName        string `json:"domain_name" form:"domain_name" query:"domain_name"`
	TenantID          string `json:"tenant_id" form:"tenant_id" query:"tenant_id"`
	VMRootPassword    string `json:"vm_root_password" form:"vm_root_password" query:"vm_root_password"`
	MysqlRootPassword string `json:"mysql_root_password" form:"mysql_root_password" query:"mysql_root_password"`
	NetworkID         string `json:"net_id" form:"net_id" query:"net_id"`
	NetworkName       string `json:"net_name" form:"net_name" query:"net_name"`
	FlavorID          string `json:"flavor_id" form:"flavor_id" query:"flavor_id"`
	ImageID           string `json:"image_id" form:"image_id" query:"image_id"`
}
