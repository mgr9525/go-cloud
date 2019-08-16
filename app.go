package gocloud

import (
	consulapi "github.com/hashicorp/consul/api"
	"github.com/xormplus/xorm"
	"gopkg.in/macaron.v1"
	"gopkg.in/mgo.v2"
	"log"
)

var (
	Db     *xorm.Engine
	Web    *macaron.Macaron
	Consul *consulapi.Client
	Mongo  *mgo.Session
	Logger *log.Logger
)
