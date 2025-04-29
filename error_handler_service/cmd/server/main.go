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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"project_chimera/error_handle_service/config"
	"project_chimera/error_handle_service/config/consul"
	"project_chimera/error_handle_service/config/db"
	"project_chimera/error_handle_service/config/rabbitmq"
	"project_chimera/error_handle_service/internal/actuators"
	"project_chimera/error_handle_service/internal/dump"
	customlogger "project_chimera/error_handle_service/pkg/logger"
	"syscall"
	"time"

	"github.com/getkin/kin-openapi/openapi2"
	"github.com/getkin/kin-openapi/openapi2conv"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title Error Handle Service API
// @version 1.0.0
// @description This API provides endpoints for error handler service which is a part of the project chimera.

// @contact.name Naveen R
// @contact.url https://naveen2070.github.io/portfolio
// @contact.email naveenrameshcud@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	config.LoadConfig()

	// Start background log archiver
	go func() {
		for {
			customlogger.LogInfo("Running scheduled log archiver...")
			customlogger.ArchiveOldLogs()
			time.Sleep(24 * time.Hour)
		}
	}()

	app := fiber.New()

	// Logger setup
	app.Use(logger.New(customlogger.InitLogger()))

	// MongoDB setup
	db.ConnectDB()

	collection := db.GetCollection("chimera_flora", "error_dump")

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
		customlogger.LogFatal(fmt.Sprintf("Failed to initialize RabbitMQ consumer:\n%s", err.Error()))
	}

	// Start consuming messages
	log.Println("Starting RabbitMQ consumer...")
	if err := dump.InitFloraDumpService(consumer, collection); err != nil {
		customlogger.LogFatal(fmt.Sprintf("Failed to start flora dump service:\n%s", err.Error()))
	}

	// Graceful shutdown handling
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		<-shutdownChan
		customlogger.LogInfo("Received shutdown signal. Cleaning up...")

		// Deregister the service from Consul
		if err := consul.DeregisterFromConsul(); err != nil {
			customlogger.LogError(fmt.Sprintf("Failed to deregister from Consul:\n%s", err.Error()))
		} else {
			customlogger.LogInfo("Service deregistered from Consul")
		}

		// Close RabbitMQ Consumer
		customlogger.LogInfo("Closing RabbitMQ Consumer...")
		consumer.Close()

		// Shutdown Fiber server
		customlogger.LogInfo("Shutting down Fiber server...")
		if err := app.Shutdown(); err != nil {
			customlogger.LogError(fmt.Sprintf("Failed to shutdown Fiber server:\n%s", err.Error()))
		}

		customlogger.LogInfo("Graceful shutdown completed.")
		os.Exit(0)
	}()

	// set up cross-origin resource sharing (CORS) middleware
	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
			AllowHeaders: "Content-Type, Authorization",
		},
	))

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
				URL:         "http://localhost:8080/error-handler/",
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

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Welcome to Fiber with MongoDB!"})
	})

	// Set up route groups
	actuatorGroup := app.Group("/actuator")

	// Register routes
	actuators.ActuatorRouter(actuatorGroup)

	err = app.Listen(":" + config.Env.AppPort)
	if err != nil {
		customlogger.LogFatal(fmt.Sprintf("Failed to start server:\n%s", err.Error()))
	}
}
