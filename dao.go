package gocloud

import (
	"github.com/xormplus/xorm"
)

type Dao struct {
	db       **xorm.Engine
	tempName string
}

type Page struct {
	Page  int64
	Size  int64
	Pages int64
	Total int64
	Data  interface{}
}

func NewDao(d **xorm.Engine, tmpName string) *Dao {
	e := new(Dao)
	e.db = d
	e.tempName = tmpName
	return e
}

func (c *Dao) SetDb(e **xorm.Engine) {
	c.db = e
}
func (c *Dao) getDb() *xorm.Engine {
	if c.db != nil {
		return *c.db
	}
	return nil
}
func (c *Dao) NewTSession(pars *map[string]interface{}) *xorm.Session {
	return c.getDb().SqlTemplateClient(c.tempName, pars)
}

func (c *Dao) FindCount(pars *map[string]interface{}) int64 {
	ret := int64(0)
	(*pars)["getCount"] = 1
	ok, err := c.NewTSession(pars).Get(&ret)
	if err != nil {
		println(err.Error())
		return 0
	}
	if ok {
		return ret
	}
	return 0
}
func (c *Dao) FindPage(ls interface{}, pars *map[string]interface{}, page int64, size interface{}) *Page {
	var pageno int64 = 1
	var sizeno int64 = 10
	var pagesno int64 = 0
	//var count=c.FindCount(pars)
	if page > 0 {
		pageno = page
	}
	if size != nil {
		switch size.(type) {
		case int:
			sizeno = int64(size.(int))
			break
		case int64:
			sizeno = size.(int64)
			break
		}
	}
	start := (pageno - 1) * sizeno
	err := c.NewTSession(pars).Limit(int(sizeno), int(start)).Find(ls)
	if err != nil {
		println(err.Error())
		return nil
	}
	count := c.FindCount(pars)
	pagest := count / sizeno
	if count%sizeno > 0 {
		pagesno = pagest + 1
	} else {
		pagesno = pagest
	}
	return &Page{
		Page:  pageno,
		Pages: pagesno,
		Size:  sizeno,
		Total: count,
		Data:  ls,
	}
}
