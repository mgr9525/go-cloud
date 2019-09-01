package gocloud

import (
	"encoding/gob"
	"github.com/go-macaron/cache"
	"github.com/go-macaron/gzip"
	"github.com/go-macaron/pongo2"
	"gopkg.in/macaron.v1"
	"gopkg.in/yaml.v2"
	"html/template"
	"io/ioutil"
	"log"
	"time"
)

func inits() {
	gob.Register(time.Time{})
	gob.Register(map[string]interface{}{})
}
func RunApp(ymlpath string, consRt func(), consfun func() []template.FuncMap) {
	inits()
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
	} else {
		Logger = &log.Logger{}
	}
	if CloudConf.Consul.Enable && CloudConf.Consul.Id != "" && CloudConf.Consul.Name != "" {
		runConsul(host, port)
	}
	if CloudConf.Db.Enable && CloudConf.Db.Driver != "" && CloudConf.Db.Url != "" {
		runDb()
	}
	if CloudConf.Mongo.Enable && CloudConf.Mongo.Url != "" {
		runMongo()
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
	if consRt != nil {
		consRt()
	}
	Web.Run(host, port)
}

func runMids(funcMap []template.FuncMap) {
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

	Web.Use(CheckContJson)
}
