package dao

import (
	"go-cloud"
	"go-cloud/examples/models"
)

var UserDao = gocloud.Dao{
	TempName: "user.dao.stpl",
	GetModel: func() interface{} {
		return new(models.User)
	},
	GetModels: func() interface{} {
		return &[]models.User{}
	},
}

func GetUserById(id int64) *models.User {
	user := new(models.User)
	//user.Id=1
	//user.Name="test"
	//user.Nick="测试"
	has, _ := gocloud.Db.SQL("select * from user where id=?", id).Get(user)
	if has {
		return user
	}
	return nil
}

//test
func FindUser(pars *map[string]interface{}) *models.User {
	user := new(models.User)
	has, _ := gocloud.Db.SqlTemplateClient("user.dao.stpl", pars).Get(user)
	if has {
		return user
	}
	return nil
}
