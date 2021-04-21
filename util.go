package gocloud

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

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

func MidAccessAllowFun(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", c.Request.Header.Get("Origin"))
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Credentials", "true")

	//放行所有OPTIONS方法
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	// 处理请求
	c.Next()
}
