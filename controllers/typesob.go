package controllers

//构造的结构体，对应http请求参数

type MysqlCluster struct {
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
  "f_net_id":"79ef3620-2fb2-4fa4-82a7-fbbd42243b4d",
  "flavor_id":"bdd9f1f1-a665-4b31-aacc-33eb1c7c6208",
  "image_id":"2b2682e3-ab47-4703-b502-022b241c658b"
}
*/

type CDHCluster struct {
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
	SeverNum          int    `json:"server_num" form:"server_num" query:"server_num"`
	AgentNum          int    `json:"agent_num" form:"agent_num" query:"agent_num"`
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
  "server_image_id":"f568ecb8-1370-4b24-96aa-58b313d2a124",
  "agent_image_id":"c61efc43-b4c6-4c97-bd71-2ab4da856789",
  "server_num":1,
  "agent_num":3
}
*/
