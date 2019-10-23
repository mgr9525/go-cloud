package controller

import (
	gocloud "github.com/mgr9525/go-cloud"
	"gopkg.in/macaron.v1"
)

type UserController struct{}

func (e *UserController) GetPath() string {
	return "/user"
}
func (e *UserController) Routes() {
	gocloud.Web.Any("/test", e.test)
}
func (e *UserController) Mid() []macaron.Handler {
	return []macaron.Handler{gocloud.AccessAllowFun}
}

func (e *UserController) test(c *macaron.Context) {
	c.PlainText(200, []byte("test:"+c.Query("n")))
}
