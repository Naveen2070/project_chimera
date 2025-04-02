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

package main

import (
	"log"
	"os"
	"os/signal"
	"project_chimera/error_handle_service/config"
	"project_chimera/error_handle_service/config/consul"
	"project_chimera/error_handle_service/config/db"
	"project_chimera/error_handle_service/internal/actuators"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.LoadConfig()
	db.ConnectDB()

	app := fiber.New()

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to Fiber with MongoDB!"})
	})

	// Logger setup
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${ip}:${port} | ${status} | ${method} | ${path} | ${latency}\n",
	}))

	// Start Consul registration
	err := consul.RegisterWithConsul()
	if err != nil {
		log.Fatalf("Error registering service with Consul: %v", err)
	}

	// Handle graceful shutdown
	go func() {
		// Create a channel to receive system interrupt signals
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

		// Block until we receive a signal
		<-signalChan

		// Deregister the service from Consul on shutdown
		err := consul.DeregisterFromConsul()
		if err != nil {
			log.Fatalf("Error deregistering service from Consul: %v", err)
		}

		log.Println("Service deregistered from Consul, shutting down gracefully.")
	}()

	// Set up route groups
	actuatorGroup := app.Group("/actuator")

	// Register routes
	actuators.ActuatorRouter(actuatorGroup)

	log.Fatal(app.Listen(":" + config.Env.AppPort))
}
