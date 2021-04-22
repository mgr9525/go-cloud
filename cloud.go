package gocloud

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	Web   *gin.Engine
	Cache *bolt.DB

	AppName string
)

func init() {
	AppName = "gocloud"
	Web = gin.Default()
	Web.FuncMap["AppName"] = func() string {
		return AppName
	}
	Web.FuncMap["GocloudTitle"] = func(tit, tits interface{}) string {
		if tits != nil {
			tms := fmt.Sprintf("%v", tits)
			if tms != "" {
				return tms
			}
		}
		if tit != nil {
			tms := fmt.Sprintf("%v", tit)
			if tms != "" {
				return fmt.Sprintf("%s-%s", tms, AppName)
			}
		}
		return AppName
	}
	Web.FuncMap["MgoIdHex"] = func(id primitive.ObjectID) string {
		return id.Hex()
	}
	Web.FuncMap["Str2Html"] = func(s string) template.HTML {
		return template.HTML(s)
	}
	Web.FuncMap["ClearHtml"] = func(s string) string {
		return ClearHTML(s)
	}

	Web.FuncMap["FmtDate"] = func(t time.Time) string {
		return t.Format("2006-01-02")
	}
	Web.FuncMap["FmtDateTime"] = func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	}
}

func Init(pths ...string) error {
	fls := "app.yml"
	if len(pths) > 0 && pths[0] != "" {
		fls = pths[0]
	}

	CloudConf = &cloudConfig{}
	err := ReadYamlConf(fls, CloudConf)
	if err != nil {
		return err
	}

	if CloudConf.Logger.Enable && CloudConf.Logger.Path != "" {
		runLogger()
	}
	if CloudConf.Cache.Enable && CloudConf.Cache.Adapter != "" {
		if err = runCache(); err != nil {
			return err
		}
	}

	return nil
}
func Run() error {
	if CloudConf == nil {
		if err := Init(); err != nil {
			return err
		}
	}
	host := "0.0.0.0"
	if CloudConf.Server.Host != "" {
		host = CloudConf.Server.Host
	}
	initFiles()
	return Web.Run(fmt.Sprintf("%s:%d", host, CloudConf.Server.Port))
}
func initFiles() {
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		Web.FuncMap = nil
	} else {
		tmpl := template.New("").Delims("{{", "}}").Funcs(Web.FuncMap)
		filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
			pth := "?" + strings.ReplaceAll(path, "\\", "/")
			if strings.HasPrefix(pth, "?templates/") && strings.HasSuffix(path, ".html") {
				//tmplfls = append(tmplfls, path)
				bts, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				t := tmpl.New(strings.ReplaceAll(pth, "?templates/", ""))
				t.Parse(string(bts))
			}
			return nil
		})
		Web.SetHTMLTemplate(tmpl)
	}
	if _, err := os.Stat("static"); !os.IsNotExist(err) {
		filepath.Walk("static", func(path string, info os.FileInfo, err error) error {
			pth := "?" + strings.ReplaceAll(path, "\\", "/")
			if strings.HasPrefix(pth, "?static/") {
				Web.StaticFile(strings.ReplaceAll(pth, "?static", ""), path)
			}
			return nil
		})
	}
}
