package gocloud

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
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

func (c *DBHelper) FindPage(ses *xorm.Session, ls interface{}, page int64, size ...int64) (*Page, error) {
	session := c.NewSession()
	count, err := c.findCount(session.Where(ses.Conds()), ls)
	session.Close()
	if err != nil {
		return nil, err
	}
	return c.findPages(ses, ls, count, page, size...)
}
func (c *DBHelper) findPages(ses *xorm.Session, ls interface{}, count, page int64, size ...int64) (*Page, error) {
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
func (c *DBHelper) FindPages(gen *PageGen, ls interface{}, page int64, size ...int64) (*Page, error) {
	var count int64
	counts := "count(*)"
	if gen.CountCols != "" {
		counts = fmt.Sprintf("count(%s)", gen.CountCols)
	}
	sqls := strings.Replace(gen.SQL, "{{select}}", counts, 1)
	sqls = strings.Replace(sqls, "{{limit}}", "", 1)
	_, err := c.GetDB().SQL(sqls, gen.Args...).Get(&count)
	if err != nil {
		return nil, err
	}

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

	starts := ""
	if start > 0 {
		starts = fmt.Sprintf("%d,", start)
	}
	ses := c.NewSession()
	defer ses.Close()
	sqls = strings.Replace(gen.SQL, "{{select}}", gen.FindCols, 1)
	if strings.Contains(sqls, "{{limit}}") {
		sqls = strings.Replace(sqls, "{{limit}}", fmt.Sprintf("LIMIT %s%d", starts, sizeno), 1)
	} else {
		sqls += fmt.Sprintf("\nLIMIT %s%d", starts, sizeno)
	}
	err = ses.SQL(sqls, gen.Args...).Find(ls)
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

// use {{select}},{{limit}}

type PageGen struct {
	SQL       string
	Args      []interface{}
	CountCols string
	FindCols  string
}
