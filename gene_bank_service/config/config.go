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

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ConsulHost                     string
	ConsulPort                     string
	ServiceHost                    string
	ServiceName                    string
	ServicePort                    int
	Status                         string
	Interval                       string
	Timeout                        string
	DeregisterCriticalServiceAfter string
	AppPort                        string
	RabbitMQurl                    string
}

var Env Config

func LoadConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, falling back to system environment variables")
	}

	servicePort, err := strconv.Atoi(os.Getenv("SERVICE_PORT"))
	if err != nil {
		log.Println("Error converting SERVICE_PORT to int:", err)
		return
	}

	// Map environment variables to the Config struct
	Env = Config{
		ConsulHost:                     os.Getenv("CONSUL_HOST"),
		ConsulPort:                     os.Getenv("CONSUL_PORT"),
		ServiceHost:                    os.Getenv("SERVICE_HOST"),
		ServiceName:                    os.Getenv("SERVICE_NAME"),
		ServicePort:                    servicePort,
		Status:                         os.Getenv("STATUS"),
		Interval:                       os.Getenv("INTERVAL"),
		Timeout:                        os.Getenv("TIMEOUT"),
		DeregisterCriticalServiceAfter: os.Getenv("DEREGISTER_CRITICAL_SERVICE_AFTER"),
		AppPort:                        os.Getenv("APP_PORT"),
		RabbitMQurl:                    os.Getenv("RABBITMQ_URL"),
	}

	log.Println("Configuration loaded successfully!")
}
