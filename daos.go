package gocloud

import (
	"errors"
	"reflect"
	"xorm.io/builder"
	"xorm.io/xorm"
)

type Page struct {
	Page  int64
	Size  int64
	Total int64
	Pages int64
	Data  interface{}
}

type SesFuncHandler = func(ses *xorm.Session)

func (c *DBHelper) findCount(e *xorm.Session, data interface{}) (int64, error) {
	if data == nil {
		return 0, errors.New("needs a pointer to a slice")
	}
	of := reflect.TypeOf(data)
	if of.Kind() == reflect.Ptr {
		of = of.Elem()
	}

	if of.Kind() == reflect.Slice {
		sty := of.Elem()
		if sty.Kind() == reflect.Ptr {
			sty = sty.Elem()
		}
		pv := reflect.New(sty)
		return e.Count(pv.Interface())
	}
	return 0, errors.New("GetCount err : not found any data")
}

func (c *DBHelper) FindPage(conds builder.Cond, ls interface{}, page int64, size ...int64) (*Page, error) {
	ses := c.NewSession()
	count, err := c.findCount(ses.Where(conds), ls)
	ses.Close()
	if err != nil {
		return nil, err
	}
	return c.FindPages(conds, ls, count, page, size...)
}
func (c *DBHelper) FindPages(conds builder.Cond, ls interface{}, count, page int64, size ...int64) (*Page, error) {
	var pageno int64 = 1
	var sizeno int64 = 10
	var pagesno int64 = 0
	//var count=c.FindCount(pars)
	if page > 0 {
		pageno = page
	}
	if len(size) > 0 && size[0] > 0 {
		sizeno = size[0]
	}
	start := (pageno - 1) * sizeno
	ses := c.NewSession()
	defer ses.Close()
	err := ses.Limit(int(sizeno), int(start)).Find(ls)
	if err != nil {
		return nil, err
	}
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
	}, nil
}
