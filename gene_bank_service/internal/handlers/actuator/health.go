package actuator

import (
	"project_chimera/gene_bank_service/internal/rabbitmq"
	"project_chimera/gene_bank_service/internal/services/actuator"

	"github.com/gofiber/fiber/v2"
)

// ActuatorHandler defines the interface for actuator handlers
type ActuatorHandler interface {
	Health(c *fiber.Ctx) error
}

// actuatorHandler is the concrete implementation of ActuatorHandler
type actuatorHandler struct {
	service    actuator.ActuatorService
	rmqHandler *rabbitmq.Handler
}

// NewActuatorHandler creates a new instance of actuatorHandler with a service dependency
func NewActuatorHandler(service actuator.ActuatorService, rmqHandler *rabbitmq.Handler) ActuatorHandler {
	return &actuatorHandler{service: service, rmqHandler: rmqHandler}
}

// Health handler for entire actuator
func (h *actuatorHandler) Health(c *fiber.Ctx) error {
	if h.rmqHandler.CheckQueueStatusHandler(c) != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Queue is not reachable"})
	} else if h.rmqHandler.CheckRabbitMQStatusHandler(c) != nil {
		return c.Status(500).JSON(fiber.Map{"error": "RabbitMQ is not reachable"})
	}
	return c.JSON(fiber.Map{"status": h.service.GetHealthMessage()})
}

// Health handler for the queue
func (h *actuatorHandler) QueueHealth(c *fiber.Ctx) error {
	return h.rmqHandler.CheckQueueStatusHandler(c)
}

// Health handler for RabbitMQ
func (h *actuatorHandler) RabbitMQHealth(c *fiber.Ctx) error {
	return h.rmqHandler.CheckRabbitMQStatusHandler(c)
}

// ActuatorRouter registers actuator-related routes
func ActuatorRouter(router fiber.Router, rmqHandler *rabbitmq.Handler) {
	service := actuator.NewActuatorService()
	handler := NewActuatorHandler(service, rmqHandler)

	router.Get("/health", handler.Health)
}
