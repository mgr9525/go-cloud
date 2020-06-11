package controller

import (
	gocloud "github.com/mgr9525/go-cloud"
	"gopkg.in/macaron.v1"
)

type IndexController struct{}

func (e *IndexController) GetPath() string {
	return ""
}
func (e *IndexController) Mid() []macaron.Handler {
	return []macaron.Handler{gocloud.AccessAllowFun}
}
func (e *IndexController) Routes() {
	gocloud.Web.Any("/", e.index)
}

func (IndexController) index(c *macaron.Context) {
	c.PlainText(200, []byte("hello world"))
}
