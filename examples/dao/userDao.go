package dao

import (
	"github.com/mgr9525/go-cloud"
	"github.com/mgr9525/go-cloud/examples/models"
)

func newUserDao() *userDao {
	e := new(userDao)
	e.Init("user.dao.stpl", e)
	return e
}

type userDao struct {
	gocloud.Dao
}

func (e *userDao) GetModel() interface{} {
	return new(models.User)
}
func (e *userDao) GetModels() interface{} {
	return &[]models.User{}
}

func (e *userDao) FindID(id int64) *models.User {
	one := e.Dao.FindID(id)
	if one == nil {
		return nil
	}
	return one.(*models.User)
}
