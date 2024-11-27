package services

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type ConsulService struct {
	client *api.Client
	config *api.Config
}

func NewConsulService(address string) (*ConsulService, error) {
	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)

	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &ConsulService{
		client: client,
		config: config,
	}, nil
}

func (c *ConsulService) Deregister(serviceID string) error {
	return c.client.Agent().ServiceDeregister(serviceID)
}

func (c *ConsulService) DiscoverService(serviceName string) ([]*api.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to discover service: %w", err)
	}

	return services, nil
}
