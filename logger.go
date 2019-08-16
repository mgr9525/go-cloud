package gocloud

import (
	"log"
	"os"
)

func runLogger() {
	logfl, err := os.OpenFile(CloudConf.Logger.Path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Println("Failed to open error log file:", err)
		return
	}
	Logger = log.New(logfl, CloudConf.Logger.Prefix, log.Ltime|log.Lshortfile)
}
