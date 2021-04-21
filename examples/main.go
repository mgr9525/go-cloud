package main

import (
	"github.com/mgr9525/go-cloud"
	"github.com/mgr9525/go-cloud/examples/route"
	"github.com/sirupsen/logrus"
)

func main() {
	//gocloud.Web.Any("/",routes.IndexHandler)
	gocloud.RegController(&route.IndexController{})
	if err := gocloud.Run(); err != nil {
		logrus.Errorf("RunApp err:%v", err)
	}
}
