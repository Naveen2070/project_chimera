package actuator

import (
	"project_chimera/gene_bank_service/internal/rabbitmq"
	"project_chimera/gene_bank_service/pkg/common"

	"github.com/gofiber/fiber/v2"
)

// ActuatorHandler defines the interface for actuator handlers
type ActuatorHandler interface {
	Health(c *fiber.Ctx) error
	QueueHealth(c *fiber.Ctx) error
	RabbitMQHealth(c *fiber.Ctx) error
}

// actuatorHandler is the concrete implementation of ActuatorHandler
type actuatorHandler struct {
	service     ActuatorService
	rmqHandlers []*rabbitmq.Handler // A slice to hold multiple RabbitMQ handlers
}

// NewActuatorHandler creates a new instance of actuatorHandler with service and RabbitMQ handler dependencies
func NewActuatorHandler(service ActuatorService, rmqHandlers []*rabbitmq.Handler) ActuatorHandler {
	return &actuatorHandler{service: service, rmqHandlers: rmqHandlers}
}

// Health handler for entire actuator
// @Summary Get actuator health status
// @Description Get actuator health status based on queue and RabbitMQ statuses
// @Tags Actuator
// @Produce json
// @Success 200 {object} common.SuccessResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /actuator/health [get]
func (h *actuatorHandler) Health(c *fiber.Ctx) error {
	// Iterate over the list of RabbitMQ handlers and check their statuses
	for _, rmqHandler := range h.rmqHandlers {
		if err := rmqHandler.CheckQueueStatusHandler(c); err != nil {
			return c.Status(500).JSON(common.ErrorResponse{Error: "Queue is not reachable"})
		} else if err := rmqHandler.CheckRabbitMQStatusHandler(c); err != nil {
			return c.Status(500).JSON(common.ErrorResponse{Error: "RabbitMQ is not reachable"})
		}
	}
	return c.JSON(common.SuccessResponse{Status: "Healthy"})
}

// Health handler for the queue
func (h *actuatorHandler) QueueHealth(c *fiber.Ctx) error {
	// Check the status of all RabbitMQ handlers for the queue
	for _, rmqHandler := range h.rmqHandlers {
		if err := rmqHandler.CheckQueueStatusHandler(c); err != nil {
			return c.Status(500).JSON(common.ErrorResponse{Error: "Queue is not reachable"})
		}
	}
	return c.JSON(common.SuccessResponse{Status: "Queue is reachable and healthy"})
}

// Health handler for RabbitMQ
func (h *actuatorHandler) RabbitMQHealth(c *fiber.Ctx) error {
	// Check the status of all RabbitMQ handlers
	for _, rmqHandler := range h.rmqHandlers {
		if err := rmqHandler.CheckRabbitMQStatusHandler(c); err != nil {
			return c.Status(500).JSON(common.ErrorResponse{Error: "RabbitMQ is not reachable"})
		}
	}
	return c.JSON(common.SuccessResponse{Status: "RabbitMQ is running and healthy"})
}

// ActuatorRouter registers actuator-related routes
func ActuatorRouter(router fiber.Router, rmqHandlers []*rabbitmq.Handler) {
	service := NewActuatorService()
	handler := NewActuatorHandler(service, rmqHandlers)

	router.Get("/health", handler.Health)
	router.Get("/health/queue", handler.QueueHealth)
	router.Get("/health/rabbitmq", handler.RabbitMQHealth)
}
