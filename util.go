package gocloud

import (
	"fmt"
	"gopkg.in/macaron.v1"
	"strings"
)

type ErrHandle func()

func RuisRecovers(name string, handle ErrHandle) {
	if err := recover(); err != nil {
		fmt.Print("ruisRecover(" + name + "):")
		fmt.Println(err) // 这里的err其实就是panic传入的内容，55
		if handle != nil {
			handle()
		}
	}
}

func GetNewMaps() map[string]interface{} {
	return make(map[string]interface{})
}

func ClearXSS(src string) string {
	rets := src
	rets = strings.Replace(rets, "`", "｀", -1)
	rets = strings.Replace(rets, "\\", "＼", -1)
	rets = strings.Replace(rets, "\"", "＂", -1)
	rets = strings.Replace(rets, "'", "＇", -1)
	return rets
}
func ClearHTML(src string) string {
	rets := src
	rets = strings.Replace(rets, "<", "＜", -1)
	rets = strings.Replace(rets, ">", "＞", -1)
	return rets
}

func AccessAllowFun(c *macaron.Context) {
	c.Resp.Header().Add("Access-Control-Allow-Origin", c.Req.Header.Get("Origin"))
	c.Resp.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Resp.Header().Add("Access-Control-Allow-Headers", "DNT,X-Mx-ReqToken,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization")
	c.Resp.Header().Add("Access-Control-Allow-Credentials", "true")

	if c.Req.Method == "OPTIONS" {
		c.PlainText(200, []byte("ok"))
	}
}
