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
	"project_chimera/error_handle_service/config/rabbitmq"
	"project_chimera/error_handle_service/internal/actuators"
	"project_chimera/error_handle_service/internal/dump"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.LoadConfig()

	app := fiber.New()

	// Logger setup
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${ip}:${port} | ${status} | ${method} | ${path} | ${latency}\n",
	}))

	// MongoDB setup
	db.ConnectDB()

	// Start Consul registration
	consul.RegisterWithConsul()

	// RabbitMQ setup
	queueName := "error_dump_queue"

	// Define retry parameters
	maxAttempts := 5
	retryDelay := 2 * time.Second

	// Initialize RabbitMQ Consumer with retries
	consumer, err := rabbitmq.NewConsumer(queueName, maxAttempts, retryDelay)
	if err != nil {
		log.Fatalf("Could not initialize RabbitMQ consumer: %v", err)
	}

	// Start consuming messages
	log.Println("Starting RabbitMQ consumer...")
	if err := dump.InitFloraDumpService(consumer); err != nil {
		log.Fatalf("Flora dump service failed to start: %v", err)
	}

	// Graceful shutdown handling
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-shutdownChan
		log.Println("Received shutdown signal. Cleaning up...")

		// Deregister the service from Consul
		if err := consul.DeregisterFromConsul(); err != nil {
			log.Printf("Error deregistering service from Consul: %v", err)
		} else {
			log.Println("Service deregistered from Consul successfully.")
		}

		// Close RabbitMQ Consumer
		log.Println("Closing RabbitMQ consumer...")
		consumer.Close()

		// Shutdown Fiber server
		log.Println("Shutting down Fiber server...")
		if err := app.Shutdown(); err != nil {
			log.Printf("Error shutting down Fiber server: %v", err)
		}

		log.Println("Application shutdown complete.")
		os.Exit(0)
	}()

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to Fiber with MongoDB!"})
	})

	// Set up route groups
	actuatorGroup := app.Group("/actuator")

	// Register routes
	actuators.ActuatorRouter(actuatorGroup)

	log.Fatal(app.Listen(":" + config.Env.AppPort))
}
