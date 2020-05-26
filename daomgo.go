package gocloud

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DaoMgo struct {
	db     **mgo.Session
	dbName string
	cNmae  string
}
type mongo struct {
	ses   *mgo.Session
	db    *mgo.Database
	cName string
}

func NewDaoMgo(d **mgo.Session, dname string, cname string) *DaoMgo {
	e := new(DaoMgo)
	e.db = d
	e.dbName = dname
	e.cNmae = cname
	return e
}

func (c *DaoMgo) SetDb(e **mgo.Session) {
	c.db = e
}
func (c *DaoMgo) GetSession() *mongo {
	if c.db != nil {
		rt := new(mongo)
		rt.ses = *c.db
		rt.db = rt.ses.DB(c.dbName)
		rt.cName = c.cNmae
		return rt
	}
	return nil
}
func (c *DaoMgo) NewSession() *mongo {
	if c.db != nil {
		rt := new(mongo)
		rt.ses = (*c.db).Copy()
		rt.db = rt.ses.DB(c.dbName)
		rt.cName = c.cNmae
		return rt
	}
	return nil
}

func (c *mongo) GetDB() *mgo.Database {
	return c.db
}
func (c *mongo) C() *mgo.Collection {
	return c.db.C(c.cName)
}
func (c *mongo) Close() {
	if c.db != nil {
		c.ses.Close()
		c.db = nil
		c.ses = nil
	}
}

func (c *mongo) FindCount(pars *bson.M) int {
	n, err := c.C().Find(*pars).Count()
	if err != nil {
		println(err.Error())
		return 0
	}
	return n
}
func (c *mongo) FindPage(ls interface{}, pars *bson.M, page int64, size interface{}, sorts ...string) *Page {
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
	if pars == nil {
		pars = &bson.M{}
	}
	q := c.C().Find(pars).Skip(int(start)).Limit(int(sizeno))
	if len(sorts) > 0 {
		q.Sort(sorts...)
	}
	err := q.All(ls)
	if err != nil {
		println(err.Error())
		return nil
	}
	count := int64(c.FindCount(pars))
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
