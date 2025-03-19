package rabbitmq

import (
	"github.com/gofiber/fiber/v2"
)

// Handler struct holds dependencies
type Handler struct {
	rpcClient *RPCClient
	queueName string
}

// NewHandler creates a new Handler instance
func NewQueueHandler(rpcClient *RPCClient, queueName string) *Handler {
	return &Handler{
		rpcClient: rpcClient,
		queueName: queueName,
	}
}

// SendCreateCommand handles HTTP requests and sends a command to RabbitMQ
func (h *Handler) SendCreateCommand(c *fiber.Ctx, cmd string) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid JSON"})
	}

	err := h.rpcClient.SendCreateCommand(h.queueName, cmd, data)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send command"})
	}

	return c.JSON(fiber.Map{"message": "Command sent"})
}

// CheckQueueStatusHandler checks if the queue is reachable
func (h *Handler) CheckQueueStatusHandler(c *fiber.Ctx) error {
	err := h.rpcClient.CheckQueueStatus(h.queueName)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Queue is not reachable"})
	}

	return c.JSON(fiber.Map{"message": "Queue is reachable"})
}

// CheckRabbitMQStatusHandler checks if RabbitMQ is running
func (h *Handler) CheckRabbitMQStatusHandler(c *fiber.Ctx) error {
	err := h.rpcClient.CheckRabbitMQStatus()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "RabbitMQ is not reachable"})
	}

	return c.JSON(fiber.Map{"message": "RabbitMQ is running"})
}
