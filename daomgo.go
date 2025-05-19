package gocloud

import (
	"context"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DaoMgo struct {
	dbcli  **qmgo.Client
	dbName string
	cNmae  string
}
type Mongo struct {
	ses   *qmgo.Session
	db    *qmgo.Database
	cName string
}

func NewDaoMgo(d **qmgo.Client, dname string, cname string) *DaoMgo {
	e := new(DaoMgo)
	e.dbcli = d
	e.dbName = dname
	e.cNmae = cname
	return e
}

func (c *DaoMgo) SetDbCli(e **qmgo.Client) {
	c.dbcli = e
}
func (c *DaoMgo) GetSession() *Mongo {
	if c.dbcli != nil && *c.dbcli != nil {
		rt := new(Mongo)
		//rt.ses = *c.db
		rt.db = (*c.dbcli).Database(c.dbName)
		rt.cName = c.cNmae
		return rt
	}
	return nil
}

/*func (c *DaoMgo) NewSession() *Mongo {
	if c.dbcli != nil {
		ses,_:=(*c.dbcli).Session()
		rt := new(Mongo)
		rt.ses = ses
		rt.db =ses.
		rt.cName = c.cNmae
		return rt
	}
	return nil
}*/

func (c *Mongo) GetDB() *qmgo.Database {
	return c.db
}
func (c *Mongo) C() *qmgo.Collection {
	return c.db.Collection(c.cName)
}
func (c *Mongo) Close() {
	if c.db != nil {
		//c.ses.Close()
		c.db = nil
		c.ses = nil
	}
}

func (c *Mongo) FindId(ctx context.Context, id interface{}) qmgo.QueryI {
	ids := id
	switch id.(type) {
	case string:
		if idt, err := primitive.ObjectIDFromHex(id.(string)); err == nil {
			ids = idt
		}
	}
	return c.C().Find(ctx, bson.M{"_id": ids})
}
func (c *Mongo) UpdateId(ctx context.Context, id interface{}, update interface{}) error {
	ids := id
	switch id.(type) {
	case string:
		if idt, err := primitive.ObjectIDFromHex(id.(string)); err == nil {
			ids = idt
		}
	}
	return c.C().UpdateId(ctx, ids, bson.M{"$set": update})
}
func (c *Mongo) UpdateOne(ctx context.Context, filter, update interface{}) error {
	return c.C().UpdateOne(ctx, filter, bson.M{"$set": update})
}
func (c *Mongo) UpdateAll(ctx context.Context, filter, update interface{}) (*qmgo.UpdateResult, error) {
	return c.C().UpdateAll(ctx, filter, bson.M{"$set": update})
}
func (c *Mongo) FindCount(ctx context.Context, pars bson.M) int64 {
	n, err := c.C().Find(ctx, pars).Count()
	if err != nil {
		println(err.Error())
		return 0
	}
	return n
}
func (c *Mongo) FindPage(ctx context.Context, ls interface{}, pars bson.M, page int64, size interface{}, sorts ...string) (*Page, error) {
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
		case int64:
			sizeno = size.(int64)
		}
	}
	start := (pageno - 1) * sizeno
	if pars == nil {
		pars = bson.M{}
	}
	q := c.C().Find(ctx, pars).Skip(start).Limit(sizeno)
	if len(sorts) > 0 {
		q.Sort(sorts...)
	}
	err := q.All(ls)
	if err != nil {
		return nil, err
	}
	count := int64(c.FindCount(ctx, pars))
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
