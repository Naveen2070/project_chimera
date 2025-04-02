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

package actuators

import (
	"github.com/gofiber/fiber/v2"
)

// ActuatorHandler defines the interface for actuator handlers
type ActuatorHandler interface {
	Health(c *fiber.Ctx) error
}

// actuatorHandler is the concrete implementation of ActuatorHandler
type actuatorHandler struct {
	service ActuatorService
}

// NewActuatorHandler creates a new instance of actuatorHandler with service and RabbitMQ handler dependencies
func NewActuatorHandler(service ActuatorService) ActuatorHandler {
	return &actuatorHandler{service: service}
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
	healthMessage := h.service.GetHealthMessage()
	return c.JSON(fiber.Map{"message": healthMessage})
}

// ActuatorRouter registers actuator-related routes
func ActuatorRouter(router fiber.Router) {
	service := NewActuatorService()
	handler := NewActuatorHandler(service)

	router.Get("/health", handler.Health)
}
