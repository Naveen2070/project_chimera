package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"project_chimera/gene_bank_service/internal/consul"
	"project_chimera/gene_bank_service/internal/handlers"
)

func main() {
	// Initialize the Fiber app
	app := fiber.New()

	// Register routes
	app.Get("/hello", handlers.HelloHandler)

	// Start Consul registration
	consul.RegisterWithConsul()

	// Start the Fiber server
	if err := app.Listen(":3001"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
