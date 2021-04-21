package gocloud

import (
	loglfshook "github.com/mgr9525/logrus-file-hook"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

func runLogger() {
	dir := "logs"
	if CloudConf.Logger.Path != "" {
		dir = filepath.Join(CloudConf.Logger.Path, "logs")
	}
	pmp := loglfshook.PathMap{
		logrus.InfoLevel:  filepath.Join(dir, "info.log"),
		logrus.ErrorLevel: filepath.Join(dir, "error.log"),
		logrus.DebugLevel: filepath.Join(dir, "debug.log"),
	}
	logrus.SetLevel(logrus.DebugLevel)
	var formatter logrus.Formatter = &logrus.TextFormatter{}
	if CloudConf.Logger.IsJson {
		formatter = &logrus.JSONFormatter{}
	}
	logrus.AddHook(loglfshook.NewLfsHook(pmp, formatter, CloudConf.Logger.FileSize, CloudConf.Logger.FileCount))
}
