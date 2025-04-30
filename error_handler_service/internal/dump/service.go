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
	"project_chimera/error_handle_service/config/db"
	"project_chimera/error_handle_service/internal/flora"
	"project_chimera/error_handle_service/pkg/common"
	logger "project_chimera/error_handle_service/pkg/logger"
	"project_chimera/error_handle_service/pkg/models"
	"strings"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/mongo"
)

// FloraDumpService defines the interface for order processing
type FloraDumpService interface {
	ProcessFloraDumpEvent(body []byte, deliveryTag uint64)
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
func (s *floraDumpService) ProcessFloraDumpEvent(body []byte, deliveryTag uint64) {
	var floraResp models.FloraResponse
	var errResp models.ErrorDataDTO

	err := json.Unmarshal(body, &floraResp)
	if err != nil {
		logger.LogError("Failed to parse message body (FloraResponse): " + err.Error())
		return
	}

	err = json.Unmarshal(body, &errResp)
	if err != nil {
		logger.LogError("Failed to parse message body (ErrorDataDTO): " + err.Error())
		return
	}

	switch {
	case strings.HasPrefix(floraResp.Pattern, "flora."):
		s.handleFloraEvents(floraResp, deliveryTag)
	case strings.HasPrefix(floraResp.Pattern, "user."):
		s.handleUserEvents(errResp, deliveryTag)
	default:
		logger.LogError("Unknown event type: " + floraResp.Pattern)
		s.acknowledgeMessage(deliveryTag)
	}
}

// Method to handle flora events
func (s *floraDumpService) handleFloraEvents(floraResp models.FloraResponse, deliveryTag uint64) {
	switch floraResp.Pattern {
	case "flora.created":
		logger.LogInfo("Processing flora.created event")

		res, err := flora.AutoFixFlora(floraResp)
		if err != nil {
			logger.LogError("Failed to fix flora data: " + err.Error() + " sending to error dump in db")
			s.saveFloraToDB(floraResp)
			s.acknowledgeMessage(deliveryTag)
			return
		} else {
			logger.LogInfo("Flora data fixed successfully and sending to upstream queue")
			s.sendMessageToQueueIfExists("flora_upstream_queue", res.Data.Data.Values, "add_flora")
			s.acknowledgeMessage(deliveryTag)
			return
		}
	case "flora.updated":
		logger.LogInfo("Processing flora.updated event")
		//TODO Handle flora update logic
		s.acknowledgeMessage(deliveryTag)
		return
	default:
		logger.LogError("Unhandled flora event type: " + floraResp.Pattern)
		s.acknowledgeMessage(deliveryTag)
	}
}

// Method to handle user signup event
func (s *floraDumpService) handleUserEvents(resp models.ErrorDataDTO, deliveryTag uint64) {
	switch resp.Pattern {
	case "user.create", "user.signup", "user.login",
		"user.delete", "user.softdelete",
		"user.getall", "user.getbyid",
		"user.update", "user.updatecredentials":

		logger.LogInfo("Processing " + resp.Pattern + " event")

		errorData, err := json.Marshal(resp.Data)
		if err != nil {
			logger.LogError("Failed to marshal error data for " + resp.Pattern + ": " + err.Error())
		} else {
			logger.LogError(resp.Pattern + " failed with error: " + string(errorData))
		}

		// Save only actual "error" events to the collection
		if resp.Data.Status != "Success" {
			s.saveToCustomCollection(resp, "chimera_user", "error_dump")
		}

		s.acknowledgeMessage(deliveryTag)
		return

	default:
		logger.LogError("Unhandled user event type: " + resp.Pattern)
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

// method to insert in custom mongoDB collection
func (s *floraDumpService) saveToCustomCollection(body models.ErrorDataDTO, dbName string, collectionName string) {
	collection := db.GetCollection(dbName, collectionName)
	document := common.ErrorDataToBson(body)
	_, err := collection.InsertOne(context.Background(), document)
	if err != nil {
		logger.LogError("Failed to insert data into" + dbName + "." + collectionName + ": " + err.Error())
	} else {
		logger.LogInfo("Data inserted into " + dbName + "." + collectionName + " successfully")
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

// Method to reject the RabbitMQ message
func (s *floraDumpService) rejectMessage(deliveryTag uint64) {
	err := s.channel.Reject(deliveryTag, true)
	if err != nil {
		logger.LogError("Failed to reject message: " + err.Error())
	} else {
		logger.LogInfo("Message rejected successfully")
	}
}

// Method to publish a message to a queue if it exists
func (s *floraDumpService) sendMessageToQueueIfExists(queueName string, data models.FloraData, pattern string) {
	message := map[string]interface{}{
		"pattern": map[string]string{
			"cmd": pattern,
		},
		"data": data,
	}
	// Marshal FloraData to JSON
	messageBody, err := json.Marshal(message)
	if err != nil {
		logger.LogError("Failed to marshal FloraData to JSON: " + err.Error())
		return
	}

	// Try declaring the queue passively (it must already exist)
	_, err = s.channel.QueueDeclarePassive(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		logger.LogError("Queue does not exist or could not be declared passively: " + err.Error())
		return
	}

	// Publish the message to the queue
	err = s.channel.Publish(
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        messageBody,
		},
	)

	if err != nil {
		logger.LogError("Failed to publish message to queue: " + err.Error())
	} else {
		logger.LogInfo("Message published to queue: " + queueName)
	}
}
