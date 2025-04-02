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
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer represents a RabbitMQ consumer
type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
}

// Singleton consumer instance
var instance *Consumer
var once sync.Once

// NewConsumer initializes the singleton RabbitMQ consumer
func NewConsumer(queueName string, maxAttempts int, retryDelay time.Duration) (*Consumer, error) {
	var err error

	once.Do(func() {
		rabbitMQURL := os.Getenv("RABBITMQ_URL")
		if rabbitMQURL == "" {
			err = fmt.Errorf("RABBITMQ_URL environment variable is not set")
			return
		}

		for i := 1; i <= maxAttempts; i++ {
			conn, connErr := amqp.Dial(rabbitMQURL)
			if connErr != nil {
				err = connErr
				log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i, maxAttempts, err)
				time.Sleep(retryDelay)
				continue
			}

			ch, chErr := conn.Channel()
			if chErr != nil {
				err = chErr
				log.Printf("Failed to open channel (attempt %d/%d): %v", i, maxAttempts, err)
				conn.Close()
				time.Sleep(retryDelay)
				continue
			}

			_, queueErr := ch.QueueDeclare(
				queueName, true, false, false, false, nil,
			)
			if queueErr != nil {
				err = queueErr
				log.Printf("Failed to declare queue (attempt %d/%d): %v", i, maxAttempts, err)
				ch.Close()
				conn.Close()
				time.Sleep(retryDelay)
				continue
			}

			// Successfully connected
			log.Printf("Connected to RabbitMQ on attempt %d/%d", i, maxAttempts)
			instance = &Consumer{conn: conn, channel: ch, queue: queueName}
			err = nil
			return
		}
	})

	return instance, err
}

// Consume starts consuming messages, but processing is handled by the service
func (c *Consumer) Consume(handler func(map[string]interface{})) error {
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

	log.Println("RabbitMQ consumer started...")

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
