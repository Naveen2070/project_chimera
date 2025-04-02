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

package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer represents a RabbitMQ consumer
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

// NewConsumer initializes a new RabbitMQ consumer with retries
func NewConsumer(queueName string, maxAttempts int, retryDelay time.Duration) (*Consumer, error) {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL environment variable is not set")
	}

	var conn *amqp.Connection
	var ch *amqp.Channel
	var err error

	for i := 0; i < maxAttempts; i++ {
		conn, err = amqp.Dial(rabbitMQURL)
		if err == nil {
			ch, err = conn.Channel()
			if err == nil {
				_, err = ch.QueueDeclare(
					queueName,
					true,  // durable
					false, // auto-delete
					false, // exclusive
					false, // no-wait
					nil,   // arguments
				)
				if err == nil {
					log.Printf("Connected to RabbitMQ (attempt %d)", i+1)
					return &Consumer{conn: conn, channel: ch, queue: queueName}, nil
				}
			}
		}

		log.Printf("Failed to connect to RabbitMQ (attempt %d): %v", i+1, err)
		time.Sleep(retryDelay)
	}

	// Return error after max attempts
	return nil, fmt.Errorf("could not connect to RabbitMQ after %d attempts: %w", maxAttempts, err)
}

// Consume listens for messages and dispatches them
func (c *Consumer) Consume(handler func(body map[string]interface{})) error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var body map[string]interface{}
			if err := json.Unmarshal(msg.Body, &body); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}
			handler(body)
		}
	}()

	return nil
}

// Close cleans up RabbitMQ resources
func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
