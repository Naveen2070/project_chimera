//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.

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
	config.Address = Config.Env.ConsulHost + ":" + Config.Env.ConsulPort
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	// Service registration
	registration := &api.AgentServiceRegistration{
		Name:    Config.Env.ServiceName,
		Address: Config.Env.ServiceHost,
		Port:    Config.Env.ServicePort,
		Tags:    []string{"fiber", "gene bank", "golang", "rabbitmq"},
		Check: &api.AgentServiceCheck{
			Status:                         Config.Env.Status,
			HTTP:                           "http://host.docker.internal:" + strconv.Itoa(Config.Env.ServicePort) + "/actuator/health",
			Interval:                       Config.Env.Interval,
			Timeout:                        Config.Env.Timeout,
			DeregisterCriticalServiceAfter: Config.Env.DeregisterCriticalServiceAfter,
		},
	}

	// Register the service
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("Error registering service with Consul: %v", err)
		return err
	}

	log.Printf("Service %s registered with Consul", Config.Env.ServiceName)
	return nil
}

// DeregisterFromConsul deregisters the service from Consul
func DeregisterFromConsul() error {
	// Setup Consul client
	config := api.DefaultConfig()
	config.Address = Config.Env.ConsulHost + ":" + Config.Env.ConsulPort
	client, err := api.NewClient(config)
	if err != nil {
		log.Printf("Error creating Consul client: %v", err)
		return err
	}

	// Deregister the service
	err = client.Agent().ServiceDeregister(Config.Env.ServiceName)
	if err != nil {
		log.Printf("Error deregistering service from Consul: %v", err)
		return err
	}

	log.Printf("Service %s deregistered from Consul", Config.Env.ServiceName)
	return nil
}
