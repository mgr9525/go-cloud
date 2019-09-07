package dao

import (
	"examples/models"
	"github.com/mgr9525/go-cloud"
)

func newUserDao() *userDao {
	e := new(userDao)
	e.TempName = "user.dao.stpl"
	e.GetModel = func() interface{} {
		return new(models.User)
	}
	e.GetModels = func() interface{} {
		return &[]models.User{}
	}
	return e
}

type userDao struct {
	gocloud.Dao
}

func (e *userDao) FindID(id int64) *models.User {
	one := e.Dao.FindID(id)
	if one == nil {
		return nil
	}
	return one.(*models.User)
}
