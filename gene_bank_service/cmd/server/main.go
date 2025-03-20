package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"project_chimera/gene_bank_service/internal/actuator"
	"project_chimera/gene_bank_service/internal/consul"
	"project_chimera/gene_bank_service/internal/rabbitmq"
)

func main() {
	// Initialize the Fiber app
	app := fiber.New()

	// Set up RabbitMQ client
	rabbitURL := "amqp://admin:naveen@2007@localhost:5672"
	queueName := "flora_upstream_queue"

	rpcClient, err := rabbitmq.NewRPCClient(rabbitURL)

	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rpcClient.Close()

	// Create queue handler
	UpstreamQueueHandler := rabbitmq.NewQueueHandler(rpcClient, queueName)

	// Set up logging middleware
	app.Use(logger.New())

	// Set up route groups
	actuatorGroup := app.Group("/actuator")

	// Register routes
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	actuator.ActuatorRouter(actuatorGroup, UpstreamQueueHandler)

	//logger setup
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
	if err := app.Listen(":5050"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
