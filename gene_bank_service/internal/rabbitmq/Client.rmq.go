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
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient handles both RPC and Ack-based messaging
type RabbitMQClient struct {
	conn       *amqp091.Connection
	channel    *amqp091.Channel
	replyQueue amqp091.Queue
}

// NewRabbitMQClient initializes the RabbitMQ connection
func NewRabbitMQClient(rabbitURL string) (*RabbitMQClient, error) {
	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	replyQueue, err := ch.QueueDeclare(
		"", false, false, true, false, nil,
	)
	if err != nil {
		conn.Close()
		ch.Close()
		return nil, err
	}

	log.Println("Connected to RabbitMQ")

	return &RabbitMQClient{
		conn:       conn,
		channel:    ch,
		replyQueue: replyQueue,
	}, nil
}

// SendRPCCommand sends a command and waits for a response (RPC pattern)
func (c *RabbitMQClient) SendRPCCommand(queueName string, cmd string, data interface{}) (interface{}, error) {
	message := map[string]interface{}{
		"pattern": map[string]string{
			"cmd": cmd,
		},
		"data": data,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	corrID := uuid.New().String()

	// Create a channel to receive responses
	responseChan := make(chan []byte)

	// Start a goroutine to listen for responses
	go func() {
		msgs, err := c.channel.Consume(
			c.replyQueue.Name, "", true, false, false, false, nil,
		)
		if err != nil {
			log.Printf("Failed to consume response: %v", err)
			close(responseChan)
			return
		}

		for msg := range msgs {
			if msg.CorrelationId == corrID {
				responseChan <- msg.Body
				break
			}
		}
		close(responseChan) // Close channel after receiving a response
	}()

	// Publish message with reply-to and correlation ID
	err = c.channel.PublishWithContext(context.Background(),
		"", queueName, false, false,
		amqp091.Publishing{
			ContentType:   "application/json",
			Body:          body,
			ReplyTo:       c.replyQueue.Name,
			CorrelationId: corrID,
		})

	if err != nil {
		log.Printf("Failed to publish RPC command: %v", err)
		return nil, err
	}

	log.Printf("Sent RPC command: %s with correlation ID: %s", string(body), corrID)

	// Wait for response or timeout
	select {
	case responseBody := <-responseChan:
		log.Println("Received RPC response:", string(responseBody))
		var response interface{}
		err := json.Unmarshal(responseBody, &response)
		if err != nil {
			return nil, err
		}
		log.Println("Parsed RPC response:", response)
		return response, nil
	case <-time.After(5 * time.Second): // Add timeout to prevent waiting indefinitely
		return nil, errors.New("timeout waiting for RPC response")
	}
}

// SendAckCommand sends a message without waiting for a response (Ack-based)
func (c *RabbitMQClient) SendAckCommand(queueName string, cmd string, data interface{}) error {
	message := map[string]interface{}{
		"pattern": map[string]string{
			"cmd": cmd,
		},
		"data": data,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	log.Println("Sent Ack-based command:", string(body))

	err = c.channel.PublishWithContext(context.Background(),
		"", queueName, false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Headers:     amqp091.Table{},
			Body:        body,
		})

	if err != nil {
		log.Printf("Failed to publish Ack command: %v", err)
		return err
	}
	return nil
}

// CheckQueueStatus checks if a queue exists
func (c *RabbitMQClient) CheckQueueStatus(queueName string) error {
	_, err := c.channel.QueueDeclarePassive(queueName, false, false, false, false, nil)
	if err != nil {
		log.Printf("Queue %s is not reachable: %v", queueName, err)
		return err
	}
	return nil
}

// CheckRabbitMQStatus checks if the connection is active
func (c *RabbitMQClient) CheckRabbitMQStatus() error {
	if c.conn.IsClosed() {
		log.Println("RabbitMQ connection is closed")
		return errors.New("RabbitMQ connection is closed")
	}
	return nil
}

// Close RabbitMQ connection
func (c *RabbitMQClient) Close() {
	c.channel.Close()
	c.conn.Close()
}
