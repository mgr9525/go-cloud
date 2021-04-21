package gocloud

import (
	"xorm.io/xorm"
)

//DBHelper db
type DBHelper struct {
	db *xorm.Engine
}

func NewDBHelper(db *xorm.Engine) *DBHelper {
	return &DBHelper{
		db: db,
	}
}

/*
GetDB : get xorm.Engine
*/
func (c *DBHelper) GetDB() *xorm.Engine {
	return c.db
}

/*
NewSession : get a new xorm.Session
*/
func (c *DBHelper) NewSession() *xorm.Session {
	return c.db.NewSession()
}
