package gocloud

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var CloudConf *cloudConfig

type cloudConfig struct {
	Server    serverConfig
	Cache     cacheConfig
	Logger    loggerConfig
	Token     tokenConfig
	Datasorce map[string]dbConfig
	Custom    Mp
}
type serverConfig struct {
	Name    string
	Host    string
	Port    int
	SysHost string `yaml:"sysHost"`

	TlsCert string `yaml:"tlsCert"`
	TlsPriv string `yaml:"tlsPriv"`
}
type cacheConfig struct {
	Enable  bool
	Adapter string
	Path    string
}
type loggerConfig struct {
	Enable    bool
	Path      string
	IsJson    bool  `yaml:"isJson"`
	FileSize  int64 `yaml:"fileSize"`
	FileCount int64 `yaml:"fileCount"`
}
type dbConfig struct {
	Enable bool
	Driver string
	Url    string
}
type tokenConfig struct {
	Enable   bool
	Httponly bool
	Secure   bool
	SameSite int `yaml:"sameSite"`
	Name     string
	Key      string
	Path     string
	Domain   string
}

func ReadYamlConf(fls string, conf interface{}) error {
	if fls == "" {
		return errors.New("param err")
	}
	data, err := ioutil.ReadFile(fls)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, conf)
	if err != nil {
		return err
	}
	return nil
}
