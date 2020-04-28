package controllers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/tuzkibug/auto-echo/base"
)

func MysqlPasswordInitial(c echo.Context) (err error) {
	p := new(MsgMysqlPasswordInitial)
	//调用echo.Context的Bind函数将请求参数和User对象进行绑定。
	if err = c.Bind(p); err != nil {
		return
	}

	mysql_ip := p.MysqlIP
	newpassword := p.Newpassword
	base.ConnectToMysql(mysql_ip, newpassword)

	return c.String(http.StatusOK, p.Newpassword)
}
