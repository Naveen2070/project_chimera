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
	"log"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// FloraDumpService defines the interface for order processing
type FloraDumpService interface {
	ProcessOrderEvent(body map[string]interface{}, deliveryTag uint64)
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
func (s *floraDumpService) ProcessOrderEvent(body map[string]interface{}, deliveryTag uint64) {
	eventType, ok := body["pattern"].(string)
	if !ok {
		log.Println("Invalid message format: missing 'pattern' field")
		return
	}

	switch eventType {
	case "flora.created":
		log.Printf("Processing flora.save event: %+v", body)
		s.saveFloraToDB(body) // Save the flora data to MongoDB
		s.acknowledgeMessage(deliveryTag)
	case "flora.updated":
		log.Printf("Processing flora.update event: %+v", body)
		//TODO Handle flora update logic
		s.acknowledgeMessage(deliveryTag)
	}
}

// Method to insert flora data into MongoDB
func (s *floraDumpService) saveFloraToDB(body map[string]interface{}) {
	// You can modify the data structure as per your MongoDB schema
	document := bson.D{
		{Key: "pattern", Value: body["pattern"]},
		{Key: "name", Value: body["name"]},
		{Key: "createdAt", Value: body["createdAt"]},
	}

	// Insert the document into the collection
	_, err := s.collection.InsertOne(context.Background(), document)
	if err != nil {
		log.Printf("Failed to insert flora data into MongoDB: %v", err)
	} else {
		log.Println("Flora data inserted into MongoDB successfully")
	}
}

// Method to acknowledge the RabbitMQ message
func (s *floraDumpService) acknowledgeMessage(deliveryTag uint64) {
	err := s.channel.Ack(deliveryTag, false)
	if err != nil {
		log.Printf("Failed to acknowledge message: %v", err)
	} else {
		log.Println("Message acknowledged successfully")
	}
}
