package controllers

type MsgMysqlCreate struct {
	Username   string `json:"username" form:"username" query:"username"`
	Password   string `json:"password" form:"password" query:"password"`
	DomainName string `json:"domain_name" form:"domain_name" query:"domain_name"`
	TenantID   string `json:"tenant_id" form:"tenant_id" query:"tenant_id"`
	MysqlName  string `json:"mysql_name" form:"mysql_name" query:"mysql_name"`
}

/*
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
	Newpassword string `json:"newpassword" form:"newpassword" query:"newpassword"`
}
