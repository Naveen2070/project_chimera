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
	"os"
	"sync"
	"time"

	logger "project_chimera/error_handle_service/pkg/logger"

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
				logger.LogError(fmt.Sprintf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i, maxAttempts, err))
				time.Sleep(retryDelay)
				continue
			}

			ch, chErr := conn.Channel()
			if chErr != nil {
				err = chErr
				logger.LogError(fmt.Sprintf("Failed to open channel (attempt %d/%d): %v", i, maxAttempts, err))
				conn.Close()
				time.Sleep(retryDelay)
				continue
			}

			// Declare the queue if not already declared
			_, queueErr := ch.QueueDeclare(
				queueName, true, false, false, false, nil,
			)
			if queueErr != nil {
				err = queueErr
				logger.LogError(fmt.Sprintf("Failed to declare queue (attempt %d/%d): %v", i, maxAttempts, err))
				ch.Close()
				conn.Close()
				time.Sleep(retryDelay)
				continue
			}

			// Successfully connected
			logger.LogInfo(fmt.Sprintf("Connected to RabbitMQ on attempt %d/%d", i, maxAttempts))
			instance = &Consumer{conn: conn, channel: ch, queue: queueName}
			err = nil
			return
		}
	})

	return instance, err
}

// Consume starts consuming messages, but processing is handled by the service
func (c *Consumer) Consume(handler func(body []byte, deliveryTag uint64)) error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",    // consumer tag
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}

	logger.LogInfo("RabbitMQ consumer started...")

	go func() {
		for msg := range msgs {
			// Pass the body and delivery tag to the handler
			handler(msg.Body, msg.DeliveryTag)
		}
	}()

	return nil
}

// Get the channel
func (c *Consumer) GetChannel() *amqp.Channel {
	return c.channel
}

// Get the connection
func (c *Consumer) GetConnection() *amqp.Connection {
	return c.conn
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

// SendMessage sends a message to a RabbitMQ queue, allowing dynamic queue names
func (c *Consumer) SendMessage(queueName string, message map[string]interface{}) error {
	// If channel is not available, return an error
	if c.channel == nil {
		return fmt.Errorf("channel is not open")
	}

	// If a queue name is provided, we declare that queue dynamically
	if queueName != "" {
		_, err := c.channel.QueueDeclare(
			queueName, true, false, false, false, nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %v", queueName, err)
		}
	}

	// Marshal the message to a JSON byte array
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// Publish the message to the specified or default queue
	err = c.channel.Publish(
		"",        // exchange (empty means default)
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	logger.LogInfo(fmt.Sprintf("Message sent to queue %s", queueName))
	return nil
}
