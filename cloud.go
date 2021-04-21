package gocloud

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strings"
)

var (
	Web   *gin.Engine
	Cache *bolt.DB
)

func init() {
	Web = gin.Default()
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
	var tmplfls []string
	filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			tmplfls = append(tmplfls, path)
		}
		return nil
	})
	Web.LoadHTMLFiles(tmplfls...)
	var staticfls []string
	filepath.Walk("static", func(path string, info os.FileInfo, err error) error {
		staticfls = append(staticfls, path)
		pth := "?" + strings.ReplaceAll(path, "\\", "/")
		if strings.HasPrefix(pth, "?static/") {
			Web.StaticFile(strings.ReplaceAll(pth, "?static", ""), path)
		}
		return nil
	})
	return Web.Run(fmt.Sprintf("%s:%d", host, CloudConf.Server.Port))
}
