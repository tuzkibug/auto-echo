package controllers

import (
	"net/http"

	"github.com/labstack/echo"
)

func Servertest(c echo.Context) (err error) {
	u := new(MsgMysqlCreate)
	if err = c.Bind(u); err != nil {
		return
	}
	return c.JSON(http.StatusOK, u)
}
