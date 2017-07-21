package elasticsearch

import (
	"errors"
	"strconv"

	consul "github.com/hashicorp/consul/api"
)

type Consul struct {
	Client *consul.Client
}

func (c *Consul) Service(service, tag string) ([]string, error) {

	var endpoints []string
	addrs, _, err := c.Client.Health().Service(service, tag, true, nil)

	if len(addrs) == 0 && err == nil {
		return endpoints, errors.New("No service with the name " + service + " is registered")
	}
	if err != nil {
		return endpoints, err
	}

	for _, v := range addrs {
		endpoints = append(endpoints, "http://"+v.Service.Address+":"+strconv.Itoa(v.Service.Port))
	}
	return endpoints, nil
}

func MakeClient() (*Consul, error) {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		return nil, err
	}
	return &Consul{
		Client: client,
	}, nil
}

// func endpoints() []string {
// 	var endpoints []string
// 	config := consul.DefaultConfig()
// 	// config.Address = "localhost:8500"
// 	c, err := consul.NewClient(config)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	addrs, meta, err := c.Health().Service("global-elastichsearch-check", "search", true, nil)

// 	if len(addrs) == 0 && err == nil {
// 		log.Println("Service not found")
// 		return endpoints
// 	}
// 	if err != nil {
// 		log.Println(err)
// 		return endpoints
// 	}
// 	log.Println(meta)
// 	for _, v := range addrs {
// 		log.Printf("%#v", v.Service)
// 		log.Printf("%#v", v.Service.Address)
// 		log.Printf("%#v", v.Service.Port)
// 		endpoints = append(endpoints, "http://"+v.Service.Address+":"+strconv.Itoa(v.Service.Port))
// 	}
// 	return endpoints
// }
