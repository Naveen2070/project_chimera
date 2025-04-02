package dump

import (
	"log"
)

// OrderService defines the interface for order processing
type FloraDumpService interface {
	ProcessOrderEvent(body map[string]interface{})
}

// orderService is the concrete implementation of OrderService
type floraDumpService struct{}

// NewOrderService creates a new OrderService instance
func NewFloraDumpService() FloraDumpService {
	return &floraDumpService{}
}

// ProcessOrderEvent handles RabbitMQ messages for orders
func (s *floraDumpService) ProcessOrderEvent(body map[string]interface{}) {
	eventType, ok := body["type"].(string)
	if !ok {
		log.Println("Invalid message format: missing 'type' field")
		return
	}

	switch eventType {
	case "flora.save":
		log.Printf("Processing order.placed event: %+v", body)
		// Add logic to handle order placement
	case "order.completed":
		log.Printf("Processing order.completed event: %+v", body)
		// Add logic to handle order completion
	default:
		log.Printf("Unknown event type: %s", eventType)
	}
}
