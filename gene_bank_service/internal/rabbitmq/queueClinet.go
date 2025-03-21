//	Copyright 2025 Naveen R
//
//		Licensed under the Apache License, Version 2.0 (the "License");
//		you may not use this file except in compliance with the License.
//		You may obtain a copy of the License at
//
//		http://www.apache.org/licenses/LICENSE-2.0
//
//		Unless required by applicable law or agreed to in writing, software
//		distributed under the License is distributed on an "AS IS" BASIS,
//		WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//		See the License for the specific language governing permissions and
//		limitations under the License.

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
func (h *Handler) SendAckRequest(data map[string]interface{}, cmd string) error {

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
