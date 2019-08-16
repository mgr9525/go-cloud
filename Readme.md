# go-cloud Golang 微服务
## 新服务
```
func main() {
	/*ymlpath:=""
	if len(os.Args)>1 {
		ymlpath=os.Args[1]
	}*/

	gocloud.RunApp("test.yml", constomRoute, customFun)
}

func customFun()[]template.FuncMap {
	println("constomFun")
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
	core.Web.Any("/", routes.IndexHandler)
}
```
## 生成数据表Struct
工具使用:https://github.com/go-xorm/cmd/
```
xorm reverse mysql root:root@tcp(localhost:3306)/test?charset=utf8 %GOPATH%/src/linskruis/go-cloud/goxorm
```