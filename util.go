package gocloud

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/macaron.v1"
	"reflect"
	"strconv"
	"strings"
)

type ErrHandle func()
type ContJSON map[string]interface{}

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

func getContJson(c *macaron.Context) (cjs ContJSON, rterr error) {
	defer RuisRecovers("getContJson", func() {
		rterr = errors.New("logic error")
	})
	pars := GetNewMaps()
	contp := c.Req.Header.Get("Content-Type")
	if !strings.HasPrefix(contp, "application/json") {
		return ContJSON{}, errors.New("content not json")
	}
	bts, err := c.Req.Body().Bytes()
	if err != nil {
		return ContJSON{}, err
	}
	err = json.Unmarshal(bts, &pars)
	if err != nil {
		return ContJSON{}, err
	}
	return pars, nil
}
func CheckContJson(c *macaron.Context) {
	cont, _ := getContJson(c)
	c.Set(reflect.TypeOf(cont), reflect.ValueOf(cont))
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

func (e ContJSON) GetString(key string) string {
	if e[key] == nil {
		return ""
	}

	return fmt.Sprintf("%v", e[key])
}

func (e ContJSON) GetInt(key string) (int, error) {
	if e[key] == nil {
		return 0, errors.New("not found")
	}

	v := e[key]
	switch v.(type) {
	case int:
		return v.(int), nil
	case string:
		return strconv.Atoi(v.(string))
	case int64:
		return int(v.(int64)), nil
	case float32:
		return int(v.(float32)), nil
	case float64:
		return int(v.(float64)), nil
	}
	return 0, errors.New("not found")
}
func (e ContJSON) GetFloat(key string) (float64, error) {
	if e[key] == nil {
		return 0, errors.New("not found")
	}

	v := e[key]
	switch v.(type) {
	case int:
		return float64(v.(int)), nil
	case string:
		return strconv.ParseFloat(v.(string), 64)
	case int64:
		return float64(v.(int64)), nil
	case float32:
		return float64(v.(float32)), nil
	case float64:
		return v.(float64), nil
	}
	return 0, errors.New("not found")
}
func (e ContJSON) GetBool(key string) bool {
	if e[key] == nil {
		return false
	}

	v := e[key]
	switch v.(type) {
	case bool:
		return v.(bool)
	}
	return false
}
