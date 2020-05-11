package controllers

//构造的结构体，对应http请求参数

type MsgMysqlCluster struct {
	OpenstackIP       string `json:"op_ip" form:"op_ip" query:"op_ip"`
	Username          string `json:"username" form:"username" query:"username"`
	Password          string `json:"password" form:"password" query:"password"`
	DomainName        string `json:"domain_name" form:"domain_name" query:"domain_name"`
	TenantID          string `json:"tenant_id" form:"tenant_id" query:"tenant_id"`
	VMRootPassword    string `json:"vm_root_password" form:"vm_root_password" query:"vm_root_password"`
	MysqlRootPassword string `json:"mysql_root_password" form:"mysql_root_password" query:"mysql_root_password"`
	NetworkID         string `json:"net_id" form:"net_id" query:"net_id"`
	NetworkName       string `json:"net_name" form:"net_name" query:"net_name"`
	FloatingNetworkID string `json:"f_net_id" form:"f_net_id" query:"f_net_id"`
	FlavorID          string `json:"flavor_id" form:"flavor_id" query:"flavor_id"`
	ImageID           string `json:"image_id" form:"image_id" query:"image_id"`
}

/*
创建mysql集群时构造的请求
header:
Content-Type:application/json
body:
{
  "op_ip":"10.10.191.250",
  "username":"yunwei",
  "password":"Cs_k0lla_!23",
  "domain_name":"default",
  "tenant_id":"88434702de204a568204bd5d1c9236d0",
  "vm_root_password":"root",
  "mysql_root_password":"root",
  "net_id":"71d7fca3-0de4-4a3b-8c83-6b63874c2912",
  "net_name":"zhujj_net",
  "f_net_id":"79ef3620-2fb2-4fa4-82a7-fbbd42243b4d"
  "flavor_id":"",
  "image_id":""
}
*/

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

type MsgCDHCluster struct {
	OpenstackIP string `json:"op_ip" form:"op_ip" query:"op_ip"`
	Username    string `json:"username" form:"username" query:"username"`
	Password    string `json:"password" form:"password" query:"password"`
	DomainName  string `json:"domain_name" form:"domain_name" query:"domain_name"`
	TenantID    string `json:"tenant_id" form:"tenant_id" query:"tenant_id"`
	//VMRootPassword    string `json:"vm_root_password" form:"vm_root_password" query:"vm_root_password"`
	//MysqlRootPassword string `json:"mysql_root_password" form:"mysql_root_password" query:"mysql_root_password"`
	NetworkID         string `json:"net_id" form:"net_id" query:"net_id"`
	NetworkName       string `json:"net_name" form:"net_name" query:"net_name"`
	FloatingNetworkID string `json:"f_net_id" form:"f_net_id" query:"f_net_id"`
	ServerFlavorID    string `json:"s_flavor_id" form:"s_flavor_id" query:"s_flavor_id"`
	AgentFlavorID     string `json:"a_flavor_id" form:"a_flavor_id" query:"a_flavor_id"`
	ServerImageID     string `json:"server_image_id" form:"server_image_id" query:"server_image_id"`
	AgentImageID      string `json:"agent_image_id" form:"agent_image_id" query:"agent_image_id"`
	//SeverNum    string `json:"server_num" form:"server_num" query:"server_num"`
	//AgentNum    string `json:"agent_num" form:"agent_num" query:"agent_num"`
}

/*
创建cdh集群时构造的请求
header:
Content-Type:application/json
body:
{
  "op_ip":"10.10.191.250",
  "username":"yunwei",
  "password":"Cs_k0lla_!23",
  "domain_name":"Default",
  "tenant_id":"88434702de204a568204bd5d1c9236d0",
  "net_id":"71d7fca3-0de4-4a3b-8c83-6b63874c2912",
  "net_name":"zhujj_net",
  "f_net_id":"79ef3620-2fb2-4fa4-82a7-fbbd42243b4d",
  "s_flavor_id":"57431ed0-7f89-4c97-bb28-c2e894bcb442",
  "a_flavor_id":"4c52ae46-88a3-4946-b749-71cd40f211da",
  "server_image_id":"aafe0a95-42ca-4690-9dbe-e1924bac3941",
  "agent_image_id":"e169da6d-c5e9-4f8f-9c3f-74b0bd03e9d9"
}
*/

type Port_detail struct {
	Status       string `json:"status" form:"status" query:"status"`
	Name         string `json:"name" form:"name" query:"name"`
	AdminStateUp bool   `json:"admin_state_up" form:"admin_state_up" query:"admin_state_up"`
	NetworkId    string `json:"network_id" form:"network_id" query:"network_id"`
	DeviceOwner  string `json:"device_owner" form:"device_owner" query:"device_owner"`
	MacAddress   string `json:"mac_address" form:"mac_address" query:"mac_address"`
	DeviceId     string `json:"device_id" form:"device_id" query:"device_id"`
}

type FIP_detail struct {
	RouterId          string      `json:"router_id" form:"router_id" query:"router_id"`
	Status            string      `json:"status" form:"status" query:"status"`
	Description       string      `json:"description" form:"description" query:"description"`
	Tags              []string    `json:"tags" form:"tags" query:"tags"`
	TenantId          string      `json:"tenant_id" form:"tenant_id" query:"tenant_id"`
	CreatedAt         string      `json:"created_at" form:"created_at" query:"created_at"`
	UpdatedAt         string      `json:"updated_at" form:"updated_at" query:"updated_at"`
	FloatingNetworkId string      `json:"floating_network_id" form:"floating_network_id" query:"floating_network_id"`
	Portdetail        Port_detail `json:"port_details" form:"port_details" query:"port_details"`
	FixedIp           string      `json:"fixed_ip_address" form:"fixed_ip_address" query:"fixed_ip_address"`
	FloatingIp        string      `json:"floating_ip_address" form:"floating_ip_address" query:"floating_ip_address"`
	RevisionNum       int         `json:"revision_number" form:"revision_number" query:"revision_number"`
	ProjectId         string      `json:"project_id" form:"project_id" query:"project_id"`
	PortId            string      `json:"port_id" form:"port_id" query:"port_id"`
	Id                string      `json:"id" form:"id" query:"id"`
	QosPolicyId       string      `json:"qos_policy_id" form:"qos_policy_id" query:"qos_policy_id"`
}

type FIP struct {
	FloatingIp FIP_detail `json:"floatingip" form:"floatingip" query:"floatingip"`
}

/*
获取浮动IP时，[]byte转json结构体需要使用的结构
*/
