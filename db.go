package gocloud

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"gopkg.in/mgo.v2"

	//_ "github.com/mattn/go-sqlite3"
	"github.com/xormplus/xorm"
	"log"
)

func runDb() {
	path := "./templates/sql/"
	if len(CloudConf.Db.Tlpath) > 0 {
		path = CloudConf.Db.Tlpath
	}
	db, err := xorm.NewEngine(CloudConf.Db.Driver, CloudConf.Db.Url)
	if err != nil {
		log.Fatal("db client error : ", err)
		return
	}
	err = db.RegisterSqlMap(xorm.Xml(path+CloudConf.Db.Driver, ".xml"))
	if err != nil {
		println("db RegisterSqlMap error : " + err.Error())
	}
	err = db.RegisterSqlTemplate(xorm.Pongo2(path+CloudConf.Db.Driver, ".stpl"))
	if err != nil {
		println("db RegisterSqlTemplate error : " + err.Error())
	}
	Db = db
}

func runMongo() {
	session, err := mgo.Dial(CloudConf.Mongo.Url)
	if err != nil {
		log.Fatal("mongodb client error : ", err)
	}
	Mongo = session
}
