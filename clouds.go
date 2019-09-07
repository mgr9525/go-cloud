package gocloud

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"gopkg.in/macaron.v1"
	"strconv"
)

func RegCloud(path string, dao *Dao, check func(c *macaron.Context)) {
	cloud := CloudApi{path: path, dao: dao}
	Web.Group("/cloud"+path, func() {
		Web.Any("/findId", cloud.findId)
		Web.Any("/delIds", cloud.delIds)
		Web.Any("/findCount", cloud.findCount)
		Web.Any("/findOne", cloud.findOne)
		Web.Any("/findList", cloud.findList)
		Web.Any("/findPage", cloud.findPage)

		Web.Any("/insert", cloud.insert)
		Web.Any("/update", cloud.update)
	}, check)
}

type CloudApi struct {
	path string
	dao  *Dao
}

func (e *CloudApi) delIds(c *macaron.Context) {
	ids := c.QueryStrings("ids")
	if len(ids) <= 0 {
		c.PlainText(511, []byte("ids 错误"))
		return
	}
	for _, it := range ids {
		id, err := strconv.ParseInt(it, 10, 64)
		if err != nil {
			e.dao.DelId(id)
		}
	}

	c.PlainText(200, []byte("{}"))
}
func (e *CloudApi) findId(c *macaron.Context) {
	ids := c.Query("id")
	if ids == "" {
		c.PlainText(511, []byte("id 错误"))
		return
	}
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		c.PlainText(511, []byte("id 格式错误"))
		return
	}

	c.JSON(200, e.dao.FindID(id))
}
func (e *CloudApi) findCount(c *macaron.Context) {
	par := &map[string]interface{}{}
	pars := c.Query("params")
	if pars != "" {
		err := json.Unmarshal([]byte(pars), par)
		if err != nil {
			c.PlainText(511, []byte("params 错误"))
			return
		}
	}
	n := e.dao.FindCount(par)
	c.PlainText(200, []byte(strconv.FormatInt(n, 10)))
}
func (e *CloudApi) findOne(c *macaron.Context) {
	par := &map[string]interface{}{}
	pars := c.Query("params")
	if pars != "" {
		err := json.Unmarshal([]byte(pars), par)
		if err != nil {
			c.PlainText(511, []byte("params 错误"))
			return
		}
	}
	c.JSON(200, e.dao.FindOne(par))
}
func (e *CloudApi) findList(c *macaron.Context) {
	par := &map[string]interface{}{}
	pars := c.Query("params")
	if pars != "" {
		err := json.Unmarshal([]byte(pars), par)
		if err != nil {
			c.PlainText(511, []byte("params 错误"))
			return
		}
	}
	c.JSON(200, e.dao.FindList(par))
}
func (e *CloudApi) findPage(c *macaron.Context) {
	par := &map[string]interface{}{}
	pages := c.Query("page")
	sizes := c.Query("size")
	pars := c.Query("params")
	if pars != "" {
		err := json.Unmarshal([]byte(pars), par)
		if err != nil {
			c.PlainText(511, []byte("params 错误"))
			return
		}
	}
	var page int64 = 1
	var size int64 = 10
	if pages != "" {
		page, _ = strconv.ParseInt(pages, 10, 64)
	}
	if sizes != "" {
		size, _ = strconv.ParseInt(sizes, 10, 64)
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	c.JSON(200, e.dao.FindPage(par, page, size))
}
func (e *CloudApi) insert(c *macaron.Context) {
	bean := c.Query("bean")
	if bean == "" {
		c.PlainText(511, []byte("get Body 错误"))
		return
	}
	beans, err := base64.StdEncoding.DecodeString(bean)
	if err != nil {
		c.PlainText(511, []byte("get Body 错误"))
		return
	}

	if e.dao.child == nil {
		c.PlainText(511, []byte("dao.child 错误"))
		return
	}

	var n int64 = 0
	info := e.dao.child.GetModel()
	err = json.Unmarshal([]byte(beans), info)
	if err == nil {
		n, _ = e.dao.Insert(info)
	}

	c.PlainText(200, []byte(strconv.FormatInt(n, 10)))
}
func (e *CloudApi) update(c *macaron.Context) {
	ids := c.Query("id")
	ismap := c.Query("ismap")
	if ids == "" {
		c.PlainText(511, []byte("id 错误"))
		return
	}
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		c.PlainText(511, []byte("id 格式错误"))
		return
	}

	bean := c.Query("bean")
	if bean == "" {
		c.PlainText(511, []byte("get Body 错误"))
		return
	}
	beans, err := base64.StdEncoding.DecodeString(bean)
	if err != nil {
		c.PlainText(511, []byte("get Body 错误"))
		return
	}

	if e.dao.child == nil {
		c.PlainText(511, []byte("dao.child 错误"))
		return
	}

	var n int64 = 0
	if ismap == "1" {
		info := map[string]interface{}{}
		buf := bytes.NewBuffer(beans)
		enc := gob.NewDecoder(buf)
		err = enc.Decode(&info)
		if err == nil {
			n, err = e.dao.Update(info, id)
		}
	} else {
		info := e.dao.child.GetModel()
		err = json.Unmarshal([]byte(beans), info)
		if err == nil {
			n, err = e.dao.Update(info, id)
		}
	}
	if err != nil {
		println("update err:" + err.Error())
	}

	c.PlainText(200, []byte(strconv.FormatInt(n, 10)))
}
