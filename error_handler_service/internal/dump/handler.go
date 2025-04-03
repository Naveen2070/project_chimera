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

package dump

import (
	"log"
	"project_chimera/error_handle_service/config/rabbitmq"

	"go.mongodb.org/mongo-driver/mongo"
)

// OrderHandler defines the interface for handling order-related requests
type FloraDumpHandler interface {
	ConsumeMessage(body map[string]interface{}, deliveryTag uint64)
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
func (h *floraDumpHandler) ConsumeMessage(body map[string]interface{}, deliveryTag uint64) {
	h.service.ProcessOrderEvent(body, deliveryTag)
}

// InitOrderService initializes the order service consumer
func InitFloraDumpService(consumer *rabbitmq.Consumer, collection *mongo.Collection) error {
	service := NewFloraDumpService(consumer.GetChannel(), collection)
	handler := NewFloraDumpHandler(service)

	// Pass handler's method to the RabbitMQ consumer
	err := consumer.Consume(handler.ConsumeMessage)
	if err != nil {
		return err
	}

	log.Println("Flora dump consumer service started...")
	return nil
}
