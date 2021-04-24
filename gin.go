package gocloud

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
)

type ErrorRes struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Dest string `json:"dest"`
}

func NewErrorRes(code int, msg string) *ErrorRes {
	return &ErrorRes{Code: code, Msg: msg}
}

func (c *ErrorRes) SetDest(dst string) *ErrorRes {
	c.Dest = dst
	return c
}

type GinController interface {
	GetPath() string // 必须"/"开头
	Routes(g gin.IRoutes)
}

func RegController(gc GinController) {
	var gp gin.IRoutes
	if Web == nil || gc == nil {
		return
	}
	gp = Web
	if len(gc.GetPath()) > 1 {
		gp = Web.Group(gc.GetPath())
	}
	gc.Routes(gp)
}
func HandleBind(fn interface{}, igrBindErr ...bool) gin.HandlerFunc {
	fnv := reflect.ValueOf(fn)
	if fnv.Kind() != reflect.Func {
		return nil
	}
	fnt := fnv.Type()
	return func(c *gin.Context) {
		nmIn := fnt.NumIn()
		inls := make([]reflect.Value, nmIn)
		inls[0] = reflect.ValueOf(c)
		for i := 1; i < nmIn; i++ {
			argt := fnt.In(i)
			argtr := argt
			if argt.Kind() == reflect.Ptr {
				argtr = argt.Elem()
			}
			if argtr.Kind() == reflect.Struct || argtr.Kind() == reflect.Map {
				argv := reflect.New(argtr)
				err := c.Bind(argv.Interface())
				if err != nil {
					if len(igrBindErr) <= 0 || !igrBindErr[0] {
						c.String(500, fmt.Sprintf("params err[%d]:%+v", i, err))
						return
					}
				}
				if argt.Kind() == reflect.Ptr {
					inls[i] = argv
				} else {
					inls[i] = argv.Elem()
				}
			}
		}
		fnv.Call(inls)
	}
}
