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
	// Second field is meta, which we will not be using
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
