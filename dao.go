package gocloud

import (
	"errors"
	"github.com/xormplus/xorm"
	"reflect"
)

type IDao interface {
	GetModel() interface{}
	GetModels() interface{}
}
type Dao struct {
	db       *xorm.Engine
	child    IDao
	tempName string
}

type Page struct {
	List  interface{}
	Page  int64
	Size  int64
	Pages int64
	Total int64
}

func (c *Dao) Init(tmpName string, child IDao) {
	c.child = child
	c.tempName = tmpName
}
func (c *Dao) SetDb(e *xorm.Engine) {
	c.db = e
}
func (c *Dao) getDb() *xorm.Engine {
	if c.db != nil {
		return c.db
	}
	return Db
}
func (c *Dao) FindOne(pars *map[string]interface{}) interface{} {
	if c.child == nil {
		return nil
	}
	m := c.child.GetModel()
	has, err := c.getDb().SqlTemplateClient(c.tempName, pars).Get(m)
	if err != nil {
		println(err.Error())
		return nil
	}
	if has {
		return m
	}
	return nil
}

func (c *Dao) FindID(id int64) interface{} {
	if id <= 0 {
		return nil
	}
	return c.FindOne(&map[string]interface{}{"id": id})
}
func (c *Dao) DelId(id int64) (int64, error) {
	if c.child == nil {
		return 0, errors.New("child is nil")
	}
	m := c.child.GetModel()
	s := reflect.ValueOf(m).Elem()
	v := reflect.ValueOf(id)
	s.FieldByName("Id").Set(v)
	return c.getDb().Delete(m)
}
func (c *Dao) DelIds(ids []int64) int64 {
	var ln int64 = 0
	for _, v := range ids {
		n, err := c.DelId(v)
		if err == nil {
			ln = ln + n
		}
	}
	return ln
}
func (c *Dao) FindList(pars *map[string]interface{}) interface{} {
	if c.child == nil {
		return nil
	}
	m := c.child.GetModels()
	err := c.getDb().SqlTemplateClient(c.tempName, pars).Find(m)

	if err != nil {
		println(err.Error())
		return nil
	}
	return m
}

func (c *Dao) FindCount(pars *map[string]interface{}) int64 {
	if c.child == nil {
		return 0
	}
	ret := int64(0)
	(*pars)["getCount"] = 1
	ok, err := c.getDb().SqlTemplateClient(c.tempName, pars).Get(&ret)
	if err != nil {
		println(err.Error())
		return 0
	}
	if ok {
		return ret
	}
	return 0
}

func (c *Dao) FindPage(pars *map[string]interface{}, page int64, size interface{}) *Page {
	if c.child == nil {
		return nil
	}
	m := c.child.GetModels()
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
	err := c.getDb().SqlTemplateClient(c.tempName, pars).Limit(int(sizeno), int(start)).Find(m)
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
		List:  m,
		Page:  page,
		Pages: pagesno,
		Size:  sizeno,
		Total: count,
	}
}

func (c *Dao) Insert(bean interface{}) (int64, error) {
	if c.child == nil {
		return 0, errors.New("child is nil")
	}

	return c.getDb().Table(c.child.GetModel()).Insert(bean)
}
func (c *Dao) Update(bean interface{}, id interface{}) (int64, error) {
	if c.child == nil {
		return 0, errors.New("child is nil")
	}
	return c.getDb().Table(c.child.GetModel()).Where("id=?", id).Update(bean)
	//return c.getDb().Table(c.GetModel()).ID(id).Update(bean)
}
func (c *Dao) GetTable() *xorm.Session {
	if c.child == nil {
		return nil
	}
	return c.getDb().Table(c.child.GetModel())
}
