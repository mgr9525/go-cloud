package gocloud

import (
	"github.com/xormplus/xorm"
)

type Page struct {
	Page  int64
	Size  int64
	Pages int64
	Total int64
	Data  interface{}
}

func XormFindPage(ses *xorm.Session, ls interface{}, page int64, size interface{}) (*Page, error) {
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
	err := ses.Limit(int(sizeno), int(start)).Find(ls)
	if err != nil {
		return nil, err
	}
	count, err := ses.Count(ls)
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
