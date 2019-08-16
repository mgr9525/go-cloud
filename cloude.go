package gocloud

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

type CloudExec struct {
	Serv      string
	Host      string
	Path      string
	Funs      func(*url.Values)
	GetModel  func() interface{}
	GetModels func() interface{}
}

func (c *CloudExec) execHttp(path string, data *url.Values) ([]byte, error) {
	host := c.Host
	if len(host) <= 0 && len(c.Serv) > 0 && Consul != nil {
		services, err := Consul.Agent().Services()
		if err != nil {
			return nil, err
		}
		service := services[c.Serv]
		if service == nil {
			return nil, errors.New("no service")
		}
		host = fmt.Sprintf("http://%s:%d", service.Address, service.Port)
	}
	res, err := http.PostForm(host+path, *data)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bts, nil
}
func (c *CloudExec) FindId(id int64) (interface{}, error) {
	if id <= 0 {
		return nil, nil
	}
	data := &url.Values{}
	c.Funs(data)
	data.Set("id", strconv.FormatInt(id, 10))
	bts, err := c.execHttp("/cloud"+c.Path+"/findId", data)
	if err != nil {
		return nil, err
	}
	e := c.GetModel()
	err = json.Unmarshal(bts, e)
	if err != nil {
		return nil, err
	}

	return e, nil

}
func (c *CloudExec) DelIds(ids []int64) (interface{}, error) {
	data := &url.Values{}
	c.Funs(data)
	for _, it := range ids {
		data.Add("ids", strconv.FormatInt(it, 10))
	}
	bts, err := c.execHttp("/cloud"+c.Path+"/delIds", data)
	if err != nil {
		return nil, err
	}
	e := c.GetModel()
	err = json.Unmarshal(bts, e)
	if err != nil {
		return nil, err
	}

	return e, nil

}
func (c *CloudExec) FindCount(pars *map[string]interface{}) (interface{}, error) {
	data := &url.Values{}
	c.Funs(data)
	par, err := json.Marshal(pars)
	if err != nil {
		return nil, err
	}
	data.Set("params", string(par))
	bts, err := c.execHttp("/cloud"+c.Path+"/findCount", data)
	if err != nil {
		return nil, err
	}
	e := c.GetModel()
	err = json.Unmarshal(bts, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (c *CloudExec) FindOne(pars *map[string]interface{}) (interface{}, error) {
	data := &url.Values{}
	c.Funs(data)
	par, err := json.Marshal(pars)
	if err != nil {
		return nil, err
	}
	data.Set("params", string(par))
	bts, err := c.execHttp("/cloud"+c.Path+"/findOne", data)
	if err != nil {
		return nil, err
	}
	if string(bts) == "null" {
		return nil, nil
	}
	e := c.GetModel()
	err = json.Unmarshal(bts, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (c *CloudExec) FindList(pars *map[string]interface{}) (interface{}, error) {
	data := &url.Values{}
	c.Funs(data)
	par, err := json.Marshal(pars)
	if err != nil {
		return nil, err
	}
	data.Set("params", string(par))
	bts, err := c.execHttp("/cloud"+c.Path+"/findList", data)
	if err != nil {
		return nil, err
	}
	e := c.GetModels()
	err = json.Unmarshal(bts, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (c *CloudExec) FindPage(pars *map[string]interface{}, page int64, size interface{}) (*Page, error) {
	data := &url.Values{}
	c.Funs(data)
	par, err := json.Marshal(pars)
	if err != nil {
		return nil, err
	}
	data.Set("page", strconv.FormatInt(page, 10))
	if size != nil {
		data.Set("size", fmt.Sprintf("%d", size))
	}
	data.Set("params", string(par))
	bts, err := c.execHttp("/cloud"+c.Path+"/findPage", data)
	if err != nil {
		return nil, err
	}
	e := &Page{List: c.GetModels()}
	err = json.Unmarshal(bts, e)
	if err != nil {
		return nil, err
	}

	return e, nil
}
func (c *CloudExec) Insert(bean interface{}) (int64, error) {
	data := &url.Values{}
	c.Funs(data)
	par, err := json.Marshal(bean)
	if err != nil {
		return 0, err
	}
	beans := base64.StdEncoding.EncodeToString(par)
	data.Set("bean", beans)
	bts, err := c.execHttp("/cloud"+c.Path+"/insert", data)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(string(bts), 10, 64)
	if err != nil {
		return 0, err
	}

	return n, nil
}
func (c *CloudExec) Update(bean interface{}, id interface{}) (int64, error) {
	data := &url.Values{}
	c.Funs(data)

	ids := ""
	switch id.(type) {
	case int:
	case int64:
		ids = fmt.Sprintf("%d", id)
	case string:
		ids = id.(string)
	}
	ismap := false
	var objmap *map[string]interface{}
	t := reflect.ValueOf(bean)
	if t.Kind() == reflect.Map {
		ismap = true
		tmp := t.Interface().(map[string]interface{})
		objmap = &tmp
	} else if t.Kind() == reflect.Ptr {
		v := t.Elem()
		if v.Kind() == reflect.Map {
			ismap = true
			tmp := v.Interface().(map[string]interface{})
			objmap = &tmp
		}
	}

	data.Set("id", ids)
	if ismap {
		var buf bytes.Buffer
		enc := gob.NewEncoder(&buf)
		err := enc.Encode(objmap)
		if err != nil {
			return 0, err
		}
		beans := base64.StdEncoding.EncodeToString(buf.Bytes())
		data.Set("ismap", "1")
		data.Set("bean", beans)
	} else {
		par, err := json.Marshal(bean)
		if err != nil {
			return 0, err
		}
		beans := base64.StdEncoding.EncodeToString(par)
		data.Set("bean", beans)
	}
	bts, err := c.execHttp("/cloud"+c.Path+"/update", data)
	if err != nil {
		return 0, err
	}
	n, err := strconv.ParseInt(string(bts), 10, 64)
	if err != nil {
		return 0, err
	}

	return n, nil
}
