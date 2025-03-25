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

package flora

import (
	"log"
	"project_chimera/gene_bank_service/internal/rabbitmq"
	"project_chimera/gene_bank_service/pkg/common"

	"github.com/gofiber/fiber/v2"
)

// FloraHandler defines the interface for flora handlers
type FloraHandler interface {
	GetFlora(c *fiber.Ctx) error
	PostFlora(c *fiber.Ctx) error
	PutFlora(c *fiber.Ctx) error
	DeleteFlora(c *fiber.Ctx) error
}

// floraHandler is the concrete implementation of FloraHandler
type floraHandler struct {
	service FloraService
}

// NewFloraHandler creates a new instance of floraHandler with RabbitMQ handler dependencies
func NewFloraHandler(service FloraService) FloraHandler {
	return &floraHandler{service: service}
}

// GetFlora handler for retrieving flora data
// @Summary Retrieve flora data from the database
// @Tags Flora
// @Accept json
// @Produce json
// @Success 200 {object} dto.FloraResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /flora [get]
func (h *floraHandler) GetFlora(c *fiber.Ctx) error {
	res, err := h.service.GetFlora(c)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(res)
}

// PostFlora handler for adding flora data
// @Summary Add a flora data to the database
// @Tags Flora
// @Accept json
// @Produce json
// @Param flora body dto.FloraRequest true "Flora data"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /flora [post]
func (h *floraHandler) PostFlora(c *fiber.Ctx) error {
	err := h.service.PostFlora(c)
	log.Println(err)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(common.SuccessResponse{Status: "Flora submitted successfully"})
}

// PutFlora handler for updating flora data
// @Summary Update a flora data in the database
// @Tags Flora
// @Accept json
// @Produce json
// @Param flora body dto.FloraUpdateRequest true "Flora data"
// @Success 200 {object} common.SuccessResponse
// @Failure 400 {object} common.ErrorResponse
// @Failure 500 {object} common.ErrorResponse
// @Router /flora [put]
func (h *floraHandler) PutFlora(c *fiber.Ctx) error {
	err := h.service.PutFlora(c)
	log.Println(err)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(common.SuccessResponse{Status: "Flora submitted successfully"})
}

// DeleteFlora handler for deleting flora data
func (h *floraHandler) DeleteFlora(c *fiber.Ctx) error {
	return nil
}

// FloraRouter sets up the routes for flora endpoints
func FloraRouter(router fiber.Router, upStreamHandler *rabbitmq.Handler, downStreamHandler *rabbitmq.Handler) {
	service := NewFloraService(upStreamHandler, downStreamHandler)
	handler := NewFloraHandler(service)

	router.Get("/", handler.GetFlora)
	router.Post("/", handler.PostFlora)
	router.Put("/", handler.PutFlora)
	router.Delete("/", handler.DeleteFlora)
}
