package main

import (
	//"net/http"

	"github.com/labstack/echo"
	"github.com/tuzkibug/testecho/controllers"
)

func main() {

	e := echo.New()

	e.GET("/test", controllers.Servertest)

	e.POST("/createmysql", controllers.Createmysql)

	e.Logger.Fatal(e.Start(":8889"))
}
