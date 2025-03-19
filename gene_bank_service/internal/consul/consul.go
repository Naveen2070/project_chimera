package consul

import (
	"log"

	Config "project_chimera/gene_bank_service/config"

	"github.com/hashicorp/consul/api"
)

func RegisterWithConsul() {
	// Setup Consul client
	config := api.DefaultConfig()
	config.Address = Config.ConsulHost + ":" + Config.ConsulPort
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	// Service registration
	registration := &api.AgentServiceRegistration{
		ID:      Config.ServiceName,
		Name:    Config.ServiceName,
		Address: "localhost", // Replace with the actual service address
		Port:    3000,        // Your service port
		Tags:    []string{"fiber"},
	}

	// Register the service
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Error registering service with Consul: %v", err)
	}

	log.Printf("Service %s registered with Consul", Config.ServiceName)
}
