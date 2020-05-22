package main

import (
	"github.com/labstack/echo"
	"github.com/tuzkibug/auto-echo/controllers"
	"github.com/tuzkibug/auto-echo/ob"
)

//echo框架web服务启动，并做路由控制

func main() {

	e := echo.New()
	//测试用
	e.GET("/test", controllers.Servertest)
	//全流程自动化拉起mysql主备集群
	e.POST("/automysql", controllers.BuilMysqlCluster)
	//全流程自动化拉起CDH集群
	e.POST("/autocdh", controllers.BuilCDHCluster)
	e.POST("/democdh", ob.BuildCDHCluster)
	e.Logger.Fatal(e.Start(":8889"))
}
