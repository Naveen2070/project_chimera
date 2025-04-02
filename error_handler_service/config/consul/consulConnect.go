// Copyright 2025 Naveen R
//
//	Licensed under the Apache License, Version 2.0 (the "License");
//	you may not use this file except in compliance with the License.
//	You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
//	Unless required by applicable law or agreed to in writing, software
//	distributed under the License is distributed on an "AS IS" BASIS,
//	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//	See the License for the specific language governing permissions and
//	limitations under the License.

package consul

import (
	"log"
	"strconv"
	"strings"

	Config "project_chimera/error_handle_service/config"

	"github.com/google/uuid"
	"github.com/hashicorp/consul/api"
)

var customUUID = uuid.New().String()
var ServiceID = "error_handle_service" + "-" + strings.Split(customUUID, "-")[0] + "-" + strings.Split(customUUID, "-")[4]

// RegisterWithConsul registers the service with Consul
func RegisterWithConsul() {
	// Setup Consul client
	config := api.DefaultConfig()
	config.Address = Config.Env.ConsulHost + ":" + Config.Env.ConsulPort
	client, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("Error creating Consul client: %v", err)
	}

	// Service registration
	registration := &api.AgentServiceRegistration{
		ID:      ServiceID,
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
	}

	log.Printf("Service %s registered with Consul", Config.Env.ServiceName)
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
	err = client.Agent().ServiceDeregister(ServiceID)
	if err != nil {
		log.Printf("Error deregistering service from Consul: %v", err)
		return err
	}

	log.Printf("Service %s deregistered from Consul", ServiceID)
	return nil
}
