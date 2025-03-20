package rabbitmq

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// Handler struct holds dependencies
type Handler struct {
	rpcClient *RabbitMQClient
	queueName string
}

// NewHandler creates a new Handler instance
func NewQueueHandler(rpcClient *RabbitMQClient, queueName string) *Handler {
	return &Handler{
		rpcClient: rpcClient,
		queueName: queueName,
	}
}

// SendRequest handles HTTP requests and sends a RPC command to RabbitMQ
func (h *Handler) SendRequest(c *fiber.Ctx, cmd string) (interface{}, error) {
	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return nil, &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid JSON"}
	}

	response, err := h.rpcClient.SendRPCCommand(h.queueName, cmd, data)
	if err != nil {
		return nil, &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Failed to send command"}
	}

	return response, nil
}

// SendAckRequest handles HTTP requests and sends an Ack-based command to RabbitMQ
func (h *Handler) SendAckRequest(c *fiber.Ctx, cmd string) error {

	var data map[string]interface{}

	if err := c.BodyParser(&data); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid JSON"}
	}
	log.Printf("Received request for Ack-based command: %s", cmd)
	err := h.rpcClient.SendAckCommand(h.queueName, cmd, data)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Failed to send command"}
	}

	return nil
}

// CheckQueueStatusHandler checks if the queue is reachable
func (h *Handler) CheckQueueStatusHandler(c *fiber.Ctx) error {
	err := h.rpcClient.CheckQueueStatus(h.queueName)
	if err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "Queue is not reachable"}
	}

	return c.JSON(fiber.Map{"message": "Queue is reachable"})
}

// CheckRabbitMQStatusHandler checks if RabbitMQ is running
func (h *Handler) CheckRabbitMQStatusHandler(c *fiber.Ctx) error {
	err := h.rpcClient.CheckRabbitMQStatus()
	if err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: "RabbitMQ is not reachable"}
	}

	return c.JSON(fiber.Map{"message": "RabbitMQ is running"})
}
