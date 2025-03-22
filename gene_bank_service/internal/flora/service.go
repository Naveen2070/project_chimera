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
	"fmt"
	"project_chimera/gene_bank_service/internal/dto"
	"project_chimera/gene_bank_service/internal/rabbitmq"
	"project_chimera/gene_bank_service/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// FloraService defines the interface for flora services
type FloraService interface {
	GetFlora(c *fiber.Ctx) error
	PostFlora(c *fiber.Ctx) error
	PutFlora(c *fiber.Ctx) error
	DeleteFlora(c *fiber.Ctx) error
}

// FloraService is the concrete implementation of FloraService
type floraService struct {
	rmqHandlers *rabbitmq.Handler
}

func NewFloraService(rmqHandlers *rabbitmq.Handler) FloraService {
	return &floraService{rmqHandlers: rmqHandlers}
}

// GetFlora handler for retrieving flora data
func (s *floraService) GetFlora(c *fiber.Ctx) error {
	return nil
}

// PostFlora handler for adding flora data
func (s *floraService) PostFlora(c *fiber.Ctx) error {
	var payload dto.FloraRequest

	if err := c.BodyParser(&payload); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid JSON"}
	}

	// Handle image conversion to byte array
	var imageBytes []byte
	var err error

	if payload.ImageURL != "" {
		// If the image is provided via a URL, fetch it
		imageBytes, err = utils.FetchImageFromURL(payload.ImageURL)
	} else if payload.ImagePath != "" {
		// If the image is provided via a local path, read it
		imageBytes, err = utils.FetchImageFromPath(payload.ImagePath)
	} else {
		// Handle case where there is no image provided
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No image URL or path provided"}
	}

	if err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: fmt.Sprintf("Error fetching image: %v", err)}
	}

	userId := c.Get("X-Auth-UserId")

	if userId == "" {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "User ID not found in request"}
	}

	// Send Ack request
	var data map[string]interface{}
	data = map[string]interface{}{
		"CommonName":     payload.CommonName,
		"ScientificName": payload.ScientificName,
		"Image":          imageBytes,
		"Description":    payload.Description,
		"Origin":         payload.Origin,
		"OtherDetails":   payload.OtherDetails,
		"Type":           payload.Type,
		"UserId":         userId,
	}
	err = s.rmqHandlers.SendAckRequest(data, "add_flora")
	if err != nil {
		return err
	}

	return nil
}

// PutFlora handler for updating flora data
func (s *floraService) PutFlora(c *fiber.Ctx) error {
	var payload dto.FloraUpdateRequest

	if err := c.BodyParser(&payload); err != nil {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "Invalid JSON"}
	}

	// Handle image conversion to byte array
	var imageBytes []byte
	var err error

	if payload.ImageURL != "" {
		// If the image is provided via a URL, fetch it
		imageBytes, err = utils.FetchImageFromURL(payload.ImageURL)
	} else if payload.ImagePath != "" {
		// If the image is provided via a local path, read it
		imageBytes, err = utils.FetchImageFromPath(payload.ImagePath)
	} else {
		// Handle case where there is no image provided
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No image URL or path provided"}
	}

	if err != nil {
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: fmt.Sprintf("Error fetching image: %v", err)}
	}

	userId := c.Get("X-Auth-UserId")

	if userId == "" {
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "User ID not found in request"}
	}

	// Send Ack request
	data := utils.CreateFloraDataMap(payload, userId, imageBytes)
	err = s.rmqHandlers.SendAckRequest(data, "update_flora")
	if err != nil {
		return err
	}

	return nil
}

// DeleteFlora handler for deleting flora data
func (s *floraService) DeleteFlora(c *fiber.Ctx) error {
	return nil
}
