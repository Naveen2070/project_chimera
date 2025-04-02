package dump

import (
	"log"
	"project_chimera/error_handle_service/config/rabbitmq"
)

// OrderHandler defines the interface for handling order-related requests
type FloraDumpHandler interface {
	ConsumeMessage(body map[string]interface{})
}

// orderHandler is the concrete implementation of OrderHandler
type floraDumpHandler struct {
	service FloraDumpService
}

// NewOrderHandler creates a new OrderHandler with the service dependency
func NewFloraDumpHandler(service FloraDumpService) FloraDumpHandler {
	return &floraDumpHandler{service: service}
}

// ConsumeMessage processes incoming RabbitMQ messages for orders
func (h *floraDumpHandler) ConsumeMessage(body map[string]interface{}) {
	h.service.ProcessOrderEvent(body)
}

// InitOrderService initializes the order service consumer
func InitFloraDumpService(consumer *rabbitmq.Consumer) error {
	service := NewFloraDumpService()
	handler := NewFloraDumpHandler(service)

	// Pass handler's method to the RabbitMQ consumer
	err := consumer.Consume(handler.ConsumeMessage)
	if err != nil {
		return err
	}

	log.Println("Flora dump consumer service started...")
	return nil
}
