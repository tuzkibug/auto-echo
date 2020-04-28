package main

import (
	"github.com/labstack/echo"
	"github.com/tuzkibug/auto-echo/controllers"
)

func main() {

	e := echo.New()
	//测试用
	e.GET("/test", controllers.Servertest)
	//创建mysql
	e.POST("/createmysql", controllers.Createmysql)
	//获取mysql IP
	e.GET("/mysqlip", controllers.MysqlIP)
	//获取mysql mac
	e.GET("/mysqlmac", controllers.MysqlMAC)
	//配置mysql root密码
	e.POST("/mysqlrootpassword", controllers.MysqlPasswordInitial)
	//测试SSH远程登录和执行命令
	e.POST("/testssh", controllers.SSH_run_cmd)
	//向远程主机sftp上传文件
	e.POST("/sftp", controllers.UploadSSH)
	//获取用户token
	e.POST("/token", controllers.Getusertoken)

	e.Logger.Fatal(e.Start(":8889"))
}
