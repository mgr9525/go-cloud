package gocloud

import (
	"context"

	"xorm.io/xorm"
)

// DBHelper db
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
NewSession : get a new xorm.Session with auto close
*/
func (c *DBHelper) NewSession(ctxs ...context.Context) *xorm.Session {
	var ctx context.Context
	if len(ctxs) > 0 {
		ctx = ctxs[0]
	}
	if ctx == nil {
		ctx = context.Background()
	}
	return c.db.Context(ctx)
}
