package gocloud

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

var CloudConf = cloudConfig{}

type cloudConfig struct {
	Server    serverConfig
	Consul    consulConfig
	Web       webConfig
	Cache     cacheConfig
	Logger    loggerConfig
	Token     tokenConfig
	Datasorce map[string]dbConfig
	Custom    map[string]interface{}
}
type serverConfig struct {
	Name string
	Host string
	Port int
}
type consulConfig struct {
	Enable  bool
	Host    string
	Port    int
	Group   string
	Reghost string
	Tags    []string
}
type webConfig struct {
	Gzip     bool
	Template string
}
type cacheConfig struct {
	Enable   bool
	Interval int
	Adapter  string
	Configs  string
}
type loggerConfig struct {
	Enable   bool
	Path     string
	Filesize int
	Filenum  int
}
type dbConfig struct {
	Enable bool
	Driver string
	Url    string
	Tlpath string
}
type tokenConfig struct {
	Enable   bool
	Httponly bool
	Name     string
	Key      string
	Path     string
	Domain   string
}

func GetCustomConf(fls string, conf interface{}) {
	cfgs := "custom.yml"
	if len(fls) > 0 {
		cfgs = fls
	}
	data, err := ioutil.ReadFile(cfgs)
	if err != nil {
		log.Fatal("config file errs : ", err)
		return
	}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		log.Fatal("config file yaml errs : ", err)
		return
	}
}
