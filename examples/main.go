package main

import (
	"github.com/mgr9525/go-cloud"
	"github.com/mgr9525/go-cloud/examples/controller"
	"github.com/mgr9525/go-cloud/examples/routes"
	"html/template"
)

func main() {
	/*ymlpath:=""
	if len(os.Args)>1 {
		ymlpath=os.Args[1]
	}*/

	gocloud.RunApp("", constomRoute, customFun)
}

func customFun() []template.FuncMap {
	println("constomFun")
	gocloud.RegController(new(controller.UserController))
	return []template.FuncMap{map[string]interface{}{
		"AppName": func() string {
			return "GoCloud"
		},
		"AppVer": func() string {
			return "1.0.0"
		},
	}}
}
func constomRoute() {
	gocloud.Web.Any("/", routes.IndexHandler)
}
