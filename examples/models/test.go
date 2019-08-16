package models

type Test struct {
	Id   int64  `xorm:"pk autoincr BIGINT(20)"`
	Haha string `xorm:"VARCHAR(255)"`
	Hehe string `xorm:"VARCHAR(255)"`
}
