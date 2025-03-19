package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"

	"project_chimera/gene_bank_service/internal/consul"
	"project_chimera/gene_bank_service/internal/handlers"
	"project_chimera/gene_bank_service/internal/handlers/actuator"
)

func main() {
	// Initialize the Fiber app
	app := fiber.New()

	// Register routes
	app.Get("/hello", handlers.HelloHandler)
	app.Get("/actuator", actuator.HealthHandler)

	// Start Consul registration
	err := consul.RegisterWithConsul()
	if err != nil {
		log.Fatalf("Error registering service with Consul: %v", err)
	}

	// Handle graceful shutdown
	go func() {
		// Create a channel to receive system interrupt signals
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

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
	if err := app.Listen(":5050"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
