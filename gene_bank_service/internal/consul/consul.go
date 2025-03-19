package consul

import (
	"log"
	"strconv"

	Config "project_chimera/gene_bank_service/config"

	"github.com/hashicorp/consul/api"
)

// RegisterWithConsul registers the service with Consul
func RegisterWithConsul() error {
	// Setup Consul client
	config := api.DefaultConfig()
	config.Address = Config.ConsulHost + ":" + Config.ConsulPort
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	// Service registration
	registration := &api.AgentServiceRegistration{
		Name:    Config.ServiceName,
		Address: Config.ServiceHost,
		Port:    Config.ServicePort,
		Tags:    []string{"fiber", "gene bank", "golang", "rabbitmq"},
		Check: &api.AgentServiceCheck{
			Status:   "passing",
			HTTP:     "http://host.docker.internal:" + strconv.Itoa(Config.ServicePort) + "/actuator/health",
			Interval: "5s",
			Timeout:  "3s",
		},
	}

	// Register the service
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Error registering service with Consul: %v", err)
		return err
	}

	log.Printf("Service %s registered with Consul", Config.ServiceName)
	return nil
}

// DeregisterFromConsul deregisters the service from Consul
func DeregisterFromConsul() error {
	// Setup Consul client
	config := api.DefaultConfig()
	config.Address = Config.ConsulHost + ":" + Config.ConsulPort
	client, err := api.NewClient(config)
	if err != nil {
		log.Printf("Error creating Consul client: %v", err)
		return err
	}

	// Deregister the service
	err = client.Agent().ServiceDeregister(Config.ServiceName)
	if err != nil {
		log.Printf("Error deregistering service from Consul: %v", err)
		return err
	}

	log.Printf("Service %s deregistered from Consul", Config.ServiceName)
	return nil
}
