package main

import (
	"github.com/gin-gonic/gin"
	"github.com/mgr9525/go-cloud"
	"github.com/mgr9525/go-cloud/examples/route"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	gocloud.Web.Use(func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() == http.StatusNotFound && c.Writer.Size() <= 0 {
			c.String(http.StatusNotFound, "ruis not found")
		}
	})
	//gocloud.Web.Any("/",routes.IndexHandler)
	gocloud.RegController(&route.IndexController{})
	if err := gocloud.Run(); err != nil {
		logrus.Errorf("RunApp err:%v", err)
	}
}
