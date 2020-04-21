package main

import (
	//"net/http"

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
	//配置mysql root密码
	e.POST("/mysqlrootpassword", controllers.MysqlPasswordInitial)
	e.Logger.Fatal(e.Start(":8889"))
}
