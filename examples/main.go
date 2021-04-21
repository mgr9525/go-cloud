package main

import (
	"github.com/mgr9525/go-cloud"
	"github.com/mgr9525/go-cloud/examples/routes"
	"github.com/sirupsen/logrus"
)

func main() {
	/*ymlpath:=""
	if len(os.Args)>1 {
		ymlpath=os.Args[1]
	}*/

	//gocloud.Web.Any("/",routes.IndexHandler)
	gocloud.RegController(&routes.IndexController{})

	if err := gocloud.Run(); err != nil {
		logrus.Errorf("RunApp err:%v", err)
	}
}
