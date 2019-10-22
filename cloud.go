package gocloud

import (
	"github.com/go-macaron/cache"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/pongo2"
	consulapi "github.com/hashicorp/consul/api"
	"gopkg.in/macaron.v1"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
)

var (
	Web    *macaron.Macaron
	Consul *consulapi.Client
)

func RunApp(ymlpath string, consRt func(), consfun func() []template.FuncMap) {
	cfgs := "app.yml"
	if len(ymlpath) > 0 {
		cfgs = ymlpath
	}
	data, err := ioutil.ReadFile(cfgs)
	if err != nil {
		println("config file errs:" + err.Error())
		return
	}

	err = yaml.Unmarshal(data, &CloudConf)
	if err != nil {
		println("config file yaml.Unmarshal errs:" + err.Error())
		return
	}

	Web = macaron.Classic()

	host := "0.0.0.0"
	port := 4000
	if CloudConf.Server.Host != "" {
		host = CloudConf.Server.Host
	}
	if CloudConf.Server.Port > 0 {
		port = CloudConf.Server.Port
	}

	if CloudConf.Logger.Enable && CloudConf.Logger.Path != "" {
		runLogger()
	}
	if CloudConf.Consul.Enable && CloudConf.Consul.Id != "" && CloudConf.Consul.Name != "" {
		runConsul(host, port)
	}

	funcMap := []template.FuncMap{map[string]interface{}{
		"AppName": func() string {
			return "GoCloud"
		},
		"AppVer": func() string {
			return "1.0.0"
		},
	}}

	if consfun != nil {
		fmp := consfun()
		if fmp != nil {
			funcMap = fmp
		}
	}
	runMids(funcMap)
	runController()
	if consRt != nil {
		consRt()
	}
	Web.Run(host, port)
}

func runMids(funcMap []template.FuncMap) {
	Web.Use(checkContJson) //contJSON 解析
	Web.Use(macaron.Logger())
	if CloudConf.Web.Gzip {
		Web.Use(gzip.Gziper())
		Web.Use(macaron.Static("static"))
	}
	if CloudConf.Web.Template == "pongo2" {
		Web.Use(pongo2.Pongoer(pongo2.Options{
			Directory: "templates",
			//AppendDirectories:
			Extensions: []string{".tmpl", ".html"},
			//Funcs:             funcMap,
			//IndentJSON:        macaron.Env != macaron.PROD,
		}))
	} else {
		Web.Use(macaron.Renderer(macaron.RenderOptions{
			Directory: "templates",
			//AppendDirectories:
			Extensions: []string{".tmpl", ".html"},
			Funcs:      funcMap,
			IndentJSON: macaron.Env != macaron.PROD,
		}))
	}
	if CloudConf.Cache.Enable && CloudConf.Cache.Adapter != "" {
		opt := cache.Options{
			Adapter:       CloudConf.Cache.Adapter,
			AdapterConfig: CloudConf.Cache.Configs,
			Interval:      CloudConf.Cache.Interval,
		}
		if CloudConf.Cache.Adapter == "redis" {
			opt.OccupyMode = false
		}
		Web.Use(cache.Cacher(opt))
	}
}
