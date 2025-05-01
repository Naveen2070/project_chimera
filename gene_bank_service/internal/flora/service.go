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
	"log"
	"project_chimera/gene_bank_service/internal/dto"
	"project_chimera/gene_bank_service/internal/rabbitmq"
	"project_chimera/gene_bank_service/pkg/utils"
	"project_chimera/gene_bank_service/pkg/utils/helpers"

	"github.com/gofiber/fiber/v2"
)

// FloraService defines the interface for flora services
type FloraService interface {
	GetFlora(c *fiber.Ctx) (dto.FloraResponse, error)
	GetFloraById(c *fiber.Ctx) (dto.FloraResponse, error)
	PostFlora(c *fiber.Ctx) error
	PutFlora(c *fiber.Ctx) error
	DeleteFlora(c *fiber.Ctx) error
}

// FloraService is the concrete implementation of FloraService
type floraService struct {
	upStreamHandler *rabbitmq.Handler

	downStreamHandler *rabbitmq.Handler

	errorQueueHandler *rabbitmq.Handler
}

func NewFloraService(upStreamHandler *rabbitmq.Handler, downStreamHandler *rabbitmq.Handler, errorQueueHandler *rabbitmq.Handler) FloraService {
	return &floraService{upStreamHandler: upStreamHandler, downStreamHandler: downStreamHandler, errorQueueHandler: errorQueueHandler}
}

// GetFlora handler for retrieving flora data
func (s *floraService) GetFlora(c *fiber.Ctx) (dto.FloraResponse, error) {
	res, err := s.downStreamHandler.SendRequest(c, "get_all_floras", "")
	if err != nil {
		log.Printf("Error in SendRequest: %v", err)

		data := map[string]interface{}{
			"code":   500,
			"status": "Internal Server Error",
			"type":   "GET",
			"data": map[string]interface{}{
				"error": err.Error(),
			},
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.getall")
		return dto.FloraResponse{}, err
	}

	if res.Code != utils.SUCCESS {
		data := map[string]interface{}{
			"code":   res.Code,
			"status": res.Status,
			"type":   "GET",
			"data":   res.Data,
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.getall")
		return dto.FloraResponse{}, helpers.HandleRPCError(res)
	}

	floraList, err := helpers.ProcessFloraData(res.Data)
	if err != nil {
		data := map[string]interface{}{
			"code":   500,
			"status": "Internal Server Error",
			"type":   "GET",
			"data": map[string]interface{}{
				"error": err.Error(),
			},
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.getall")
		return dto.FloraResponse{}, err
	}

	return dto.FloraResponse{Flora: floraList}, nil
}

func (s *floraService) GetFloraById(c *fiber.Ctx) (dto.FloraResponse, error) {
	res, err := s.downStreamHandler.SendRequest(c, "get_flora_by_id", c.Params("id"))
	if err != nil {
		data := map[string]interface{}{
			"code":   500,
			"status": "Internal Server Error",
			"type":   "GET",
			"data": map[string]interface{}{
				"error": err.Error(),
			},
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.getbyid")
		log.Printf("Error in SendRequest: %v", err)
		return dto.FloraResponse{}, err
	}

	if res.Code != utils.SUCCESS {
		data := map[string]interface{}{
			"code":   res.Code,
			"status": res.Status,
			"type":   "GET",
			"data":   res.Data,
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.getbyid")
		return dto.FloraResponse{}, helpers.HandleRPCError(res)
	}

	floraList, err := helpers.ProcessFloraData(res.Data)
	if err != nil {
		data := map[string]interface{}{
			"code":   500,
			"status": "Internal Server Error",
			"type":   "GET",
			"data": map[string]interface{}{
				"error": err.Error(),
			},
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.getbyid")
		return dto.FloraResponse{}, err
	}

	return dto.FloraResponse{Flora: floraList}, nil
}

// PostFlora handler for adding flora data
func (s *floraService) PostFlora(c *fiber.Ctx) error {
	var payload dto.FloraRequest

	if err := c.BodyParser(&payload); err != nil {
		data := map[string]interface{}{
			"code":   400,
			"status": "Bad Request",
			"type":   "POST",
			"data": map[string]interface{}{
				"error": err.Error(),
			},
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.post")
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
	} else if payload.Image != nil {
		// If the image is provided as a byte array, use it directly
		imageBytes = payload.Image
	} else {
		data := map[string]interface{}{
			"code":   400,
			"status": "Bad Request",
			"type":   "POST",
			"data": map[string]interface{}{
				"error": "No image provided",
			},
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.post")
		// Handle case where there is no image provided
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No image URL or path provided"}
	}

	if err != nil {
		data := map[string]interface{}{
			"code":   500,
			"status": "Internal Server Error",
			"type":   "POST",
			"data": map[string]interface{}{
				"error": err.Error(),
				"url":   payload.ImageURL,
			},
		}
		s.errorQueueHandler.SendAckRequest(data, "flora.post")
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: fmt.Sprintf("Error fetching image: %v", err)}
	}

	userId := c.Get("X-Auth-UserId")

	if userId == "" {
		data := map[string]interface{}{
			"code":   400,
			"status": "Bad Request",
			"type":   "POST",
			"data": map[string]interface{}{
				"error": "User ID not found in request",
			},
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.post")
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "User ID not found in request"}
	}

	// Send Ack request
	var data map[string]interface{} = map[string]interface{}{
		"CommonName":     payload.CommonName,
		"ScientificName": payload.ScientificName,
		"Image":          imageBytes,
		"Description":    payload.Description,
		"Origin":         payload.Origin,
		"OtherDetails":   payload.OtherDetails,
		"type":           payload.Type,
		"UserId":         userId,
	}
	err = s.upStreamHandler.SendAckRequest(data, "add_flora")
	if err != nil {
		errData := map[string]interface{}{
			"code":   500,
			"status": "Internal Server Error",
			"type":   "POST",
			"data": map[string]interface{}{
				"error": err.Error(),
			},
		}

		s.errorQueueHandler.SendAckRequest(errData, "flora.post")
		return err
	}

	return nil
}

// PutFlora handler for updating flora data
func (s *floraService) PutFlora(c *fiber.Ctx) error {
	var payload dto.FloraUpdateRequest

	if err := c.BodyParser(&payload); err != nil {
		data := map[string]interface{}{
			"code":   400,
			"status": "Bad Request",
			"type":   "PUT",
			"data": map[string]interface{}{
				"error": err.Error(),
			},
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.put")
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
	} else if len(payload.Image) > 0 {
		imageBytes = payload.Image
	} else {
		// Handle case where there is no image provided
		data := map[string]interface{}{
			"code":   400,
			"status": "Bad Request",
			"type":   "PUT",
			"data": map[string]interface{}{
				"error": "No image provided",
			},
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.put")
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "No image URL or path provided"}
	}

	if err != nil {
		data := map[string]interface{}{
			"code":   500,
			"status": "Internal Server Error",
			"type":   "PUT",
			"data": map[string]interface{}{
				"error": err.Error(),
				"url":   payload.ImageURL,
			},
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.put")
		return &fiber.Error{Code: fiber.StatusInternalServerError, Message: fmt.Sprintf("Error fetching image: %v", err)}
	}

	userId := c.Get("X-Auth-UserId")

	if userId == "" {
		data := map[string]interface{}{
			"code":   400,
			"status": "Bad Request",
			"type":   "PUT",
			"data": map[string]interface{}{
				"error": "User ID not found in request",
			},
		}

		s.errorQueueHandler.SendAckRequest(data, "flora.put")
		return &fiber.Error{Code: fiber.StatusBadRequest, Message: "User ID not found in request"}
	}

	// Send Ack request
	data := utils.CreateFloraDataMap(payload, userId, imageBytes)
	err = s.upStreamHandler.SendAckRequest(data, "update_flora")
	if err != nil {
		errData := map[string]interface{}{
			"code":   500,
			"status": "Internal Server Error",
			"type":   "PUT",
			"data": map[string]interface{}{
				"error": err.Error(),
			},
		}

		s.errorQueueHandler.SendAckRequest(errData, "flora.put")
		return err
	}

	return nil
}

// DeleteFlora handler for deleting flora data
func (s *floraService) DeleteFlora(c *fiber.Ctx) error {
	return nil
}
