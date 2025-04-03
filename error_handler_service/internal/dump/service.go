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
	"context"
	"encoding/json"
	"project_chimera/error_handle_service/pkg/common"
	logger "project_chimera/error_handle_service/pkg/logger"
	"project_chimera/error_handle_service/pkg/models"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

// FloraDumpService defines the interface for order processing
type FloraDumpService interface {
	ProcessOrderEvent(body []byte, deliveryTag uint64)
}

// floraDumpService is the concrete implementation of FloraDumpService
type floraDumpService struct {
	channel    *amqp091.Channel
	collection *mongo.Collection
}

// NewFloraDumpService creates a new FloraDumpService instance with MongoDB integration
func NewFloraDumpService(channel *amqp091.Channel, collection *mongo.Collection) FloraDumpService {
	return &floraDumpService{
		channel:    channel,
		collection: collection,
	}
}

// ProcessOrderEvent handles RabbitMQ messages for orders
func (s *floraDumpService) ProcessOrderEvent(body []byte, deliveryTag uint64) {
	var floraResp models.FloraResponse

	// Parse JSON into the FloraResponse struct
	err := json.Unmarshal(body, &floraResp)
	if err != nil {
		logger.LogError("Failed to parse message body: " + err.Error())
		return
	}

	eventType := floraResp.Pattern

	switch eventType {
	case "flora.created":
		logger.LogInfo("Processing flora.created event")
		s.saveFloraToDB(floraResp)
		s.acknowledgeMessage(deliveryTag)
	case "flora.updated":
		logger.LogInfo("Processing flora.updated event")
		//TODO Handle flora update logic
		s.acknowledgeMessage(deliveryTag)
	}
}

// Method to insert flora data into MongoDB
func (s *floraDumpService) saveFloraToDB(body models.FloraResponse) {
	// You can modify the data structure as per your MongoDB schema
	document := common.FloraResponseToBson(body)

	// Insert the document into the collection
	_, err := s.collection.InsertOne(context.Background(), document)
	if err != nil {
		logger.LogError("Failed to insert flora data into MongoDB: " + err.Error())
	} else {
		logger.LogInfo("Flora data inserted into MongoDB successfully")
	}
}

// Method to acknowledge the RabbitMQ message
func (s *floraDumpService) acknowledgeMessage(deliveryTag uint64) {
	err := s.channel.Ack(deliveryTag, false)
	if err != nil {
		logger.LogError("Failed to acknowledge message: " + err.Error())
	} else {
		logger.LogInfo("Message acknowledged successfully")
	}
}
