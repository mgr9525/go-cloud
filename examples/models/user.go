package models

import (
	"time"
)

type User struct {
	Id    int64     `xorm:"pk autoincr BIGINT(20)"`
	Name  string    `xorm:"VARCHAR(100)"`
	Nick  string    `xorm:"VARCHAR(255)"`
	Times time.Time `xorm:"default 'CURRENT_TIMESTAMP' DATETIME"`
}
