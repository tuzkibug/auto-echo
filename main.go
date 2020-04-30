package main

import (
	"github.com/labstack/echo"
	"github.com/tuzkibug/auto-echo/controllers"
)

//echo框架web服务启动，并做路由控制

func main() {

	e := echo.New()
	//测试用
	e.GET("/test", controllers.Servertest)
	//全流程自动化拉起mysql主备集群
	e.POST("/automysql", controllers.BuilMysqlCluster)
	//创建mysql
	//e.POST("/createmysql", controllers.Createmysql)
	//获取mysql IP
	//e.GET("/mysqlip", controllers.MysqlIP)
	//获取mysql mac
	e.GET("/mysqlmac", controllers.MysqlMAC)
	//获取mysql port id
	e.GET("/mysqlportid", controllers.Getportid)
	//配置mysql root密码
	//e.POST("/mysqlrootpassword", controllers.MysqlPasswordInitial)
	//测试SSH远程登录和执行命令
	//e.POST("/testssh", controllers.SSH_run_cmd)
	//向远程主机sftp上传文件
	//e.POST("/sftp", controllers.UploadSSH)
	//获取用户token
	//e.POST("/token", controllers.Getusertoken)

	e.Logger.Fatal(e.Start(":8889"))
}
