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
	e.POST("/demomysql", controllers.BuildMysqlCluster)
	//全流程自动化拉起CDH集群
	e.POST("/democdh", controllers.BuildCDHCluster)

	e.Logger.Fatal(e.Start(":8889"))
}
