package gocloud

import (
	"github.com/donnie4w/go-logger/logger"
	"path/filepath"
)

func runLogger() {
	dir := filepath.Dir(CloudConf.Logger.Path)
	name := filepath.Base(CloudConf.Logger.Path)
	sz := 1
	num := 10
	if CloudConf.Logger.Filesize > 0 {
		sz = CloudConf.Logger.Filesize
	}
	if CloudConf.Logger.Filenum > 0 {
		num = CloudConf.Logger.Filenum
	}
	logger.SetRollingFile(dir, name, int32(num), int64(sz), logger.MB)
	logger.SetLevel(logger.INFO)
	//logger.Debug("logger start:",time.Now())
}
