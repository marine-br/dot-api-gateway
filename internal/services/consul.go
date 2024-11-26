package services

import (
	"fmt"
	"log"
	"net"

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

func (c *ConsulService) Register(serviceName, serviceHost string, servicePort int) error {
	ip, err := getOutboundIP()
	if err != nil {
		return fmt.Errorf("failed to get outbound IP: %w", err)
	}

	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, ip, servicePort),
		Name:    serviceName,
		Address: ip,
		Port:    servicePort,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/health", ip, servicePort),
			Interval: "10s",
			Timeout:  "5s",
		},
		Tags: []string{"api-gateway", "v1"},
	}

	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	log.Printf("Service registered in Consul: %s (%s:%d)", serviceName, ip, servicePort)
	return nil
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

func getOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")

	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String(), nil
}
