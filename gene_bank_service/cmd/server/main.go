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

package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"

	"project_chimera/gene_bank_service/config"
	"project_chimera/gene_bank_service/internal/actuator"
	"project_chimera/gene_bank_service/internal/consul"
	"project_chimera/gene_bank_service/internal/flora"
	"project_chimera/gene_bank_service/internal/rabbitmq"
)

// @title Gene Bank Service API
// @version 1.0.0
// @description This API provides endpoints for gene bank service which is a part of the project chimera.

// @contact.name Naveen R
// @contact.url https://naveen2070.github.io/portfolio
// @contact.email naveenrameshcud@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	config.LoadConfig()
	// Set up RabbitMQ client
	rabbitURL := config.Env.RabbitMQurl
	folraQueueName := "flora_upstream_queue"
	floraDownstreamQueueName := "flora_downstream_queue"

	rpcClient, err := rabbitmq.NewRabbitMQClient(rabbitURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rpcClient.Close()

	//start consumer for flora_upstream_queue
	rpcClient.StartConsumer()
	log.Println("RabbitMQ consumer started successfully!")

	// Create queue handler
	FloraUpstreamQueueHandler := rabbitmq.NewQueueHandler(rpcClient, folraQueueName)
	floraDownstreamQueueHandler := rabbitmq.NewQueueHandler(rpcClient, floraDownstreamQueueName)

	// List of RabbitMQ handlers
	var rmqHandlers = []*rabbitmq.Handler{
		FloraUpstreamQueueHandler,
		floraDownstreamQueueHandler,
	}

	// Initialize the Fiber app
	app := fiber.New()

	// set up cross-origin resource sharing (CORS) middleware
	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders: "Content-Type, Authorization",
		},
	))

	// Set up logging middleware
	app.Use(logger.New())

	// Serve Swagger JSON file (convert from Swagger 2.0 to OpenAPI 3.0)
	app.Get("swagger/v1/swagger.json", func(c *fiber.Ctx) error {
		// Read the Swagger 2.0 file
		data, err := os.ReadFile("./docs/swagger.json")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to read Swagger JSON")
		}

		// Unmarshal the JSON into an OpenAPI 2.0 structure
		var doc openapi2.T
		if err := json.Unmarshal(data, &doc); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to parse Swagger JSON")
		}

		// Convert OpenAPI 2.0 to OpenAPI 3.0
		openapi3Doc, err := openapi2conv.ToV3(&doc)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to convert to OpenAPI 3.0")
		}

		// add server info
		openapi3Doc.Servers = []*openapi3.Server{
			{
				URL:         "http://localhost:5050",
				Description: "Development Server",
			}, {
				URL:         "http://localhost:8080/gene-bank/",
				Description: "Gateway Server",
			},
		}

		// add security info
		openapi3Doc.Security = openapi3.SecurityRequirements{
			{
				"BearerAuth": []string{}, // Associates the BearerAuth with the endpoint
			},
		}
		// Define SecuritySchemes
		openapi3Doc.Components.SecuritySchemes = openapi3.SecuritySchemes{
			"BearerAuth": &openapi3.SecuritySchemeRef{
				Value: &openapi3.SecurityScheme{
					Type:         "http",
					Scheme:       "bearer",
					BearerFormat: "JWT",
				},
			},
		}

		// Serve the OpenAPI 3.0 JSON response
		return c.JSON(openapi3Doc)
	})

	// Serve Swagger UI route (Static assets)
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL: "/swagger/v1/swagger.json",
	}))

	// Set up route groups
	actuatorGroup := app.Group("/actuator")
	floraGroup := app.Group("/flora")

	// Register routes
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	actuator.ActuatorRouter(actuatorGroup, rmqHandlers)
	flora.FloraRouter(floraGroup, FloraUpstreamQueueHandler, floraDownstreamQueueHandler)

	// Logger setup
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${ip}:${port} | ${status} | ${method} | ${path} | ${latency}\n",
	}))

	// Start Consul registration
	err = consul.RegisterWithConsul()
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

	// Start the Fiber server
	if err := app.Listen(":" + config.Env.AppPort); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
