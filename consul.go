package gocloud

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"log"
)

func runConsul(host string, port int) {
	config := consulapi.DefaultConfig()
	if CloudConf.Consul.Host != "" && CloudConf.Consul.Port > 0 {
		config.Address = fmt.Sprintf("%s:%d", CloudConf.Consul.Host, CloudConf.Consul.Port)
	}
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
		return
	}
	if CloudConf.Consul.Reghost != "" {
		host = CloudConf.Consul.Reghost
	}

	//创建一个新服务。
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = CloudConf.Consul.Id
	registration.Name = CloudConf.Consul.Name
	registration.Address = host
	registration.Port = port
	registration.Tags = CloudConf.Consul.Tags
	registration.Check = &consulapi.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, registration.Port, "/consul/check"),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "30s", //check失败后30秒删除本服务
	}
	log.Println("get check.HTTP:", registration.Check)

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("register server error : ", err)
		return
	}

	Web.Get("/consul/check", func() string { return "ruisConsulCheck" })
	Consul = client
}
