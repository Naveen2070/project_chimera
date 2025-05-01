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
	"project_chimera/gene_bank_service/pkg/common"
	"project_chimera/gene_bank_service/pkg/utils"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient handles both RPC and Ack-based messaging
type RabbitMQClient struct {
	conn        *amqp091.Connection
	channel     *amqp091.Channel
	replyQueue  amqp091.Queue
	responseMap sync.Map
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

// StartConsumer listens for responses in the background
func (c *RabbitMQClient) StartConsumer() {
	go func() {
		msgs, err := c.channel.Consume(
			c.replyQueue.Name, "", false, false, false, false, nil, // Manual ACK
		)
		if err != nil {
			log.Fatalf("Failed to start consumer: %v", err)
			return
		}

		for msg := range msgs {
			if ch, ok := c.responseMap.Load(msg.CorrelationId); ok {
				responseChan := ch.(chan []byte)
				responseChan <- msg.Body                // Send response to the corresponding request
				close(responseChan)                     // Close channel after sending response
				c.responseMap.Delete(msg.CorrelationId) // Remove from map
				msg.Ack(false)                          // ✅ Manually acknowledge message
			}
		}
	}()
}

// SendRPCCommand sends an RPC request and waits for the response
func (c *RabbitMQClient) SendRPCCommand(queueName string, cmd string, data interface{}) (common.MessageResponse, error) {
	message := common.MessageRequest{
		Pattern: common.Pattern{
			Cmd: cmd,
		},
		Data: data,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return common.MessageResponse{}, err
	}

	corrID := uuid.New().String()
	responseChan := make(chan []byte, 1) // Buffered to prevent blocking

	// Store response channel in map
	c.responseMap.Store(corrID, responseChan)

	// Publish message
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
		c.responseMap.Delete(corrID) // Clean up if publish fails
		return common.MessageResponse{}, err
	}

	log.Printf("Sent RPC command: %s with correlation ID: %s", string(body), corrID)

	// Wait for response or timeout
	select {
	case responseBody := <-responseChan:
		log.Println("Received RPC response successfully")

		var response map[string]interface{}
		if err := json.Unmarshal(responseBody, &response); err != nil {
			log.Printf("Failed to parse RPC response: %v", err)
			return common.MessageResponse{}, err
		}
		log.Println("Parsed RPC response successfully")

		status, ok := response["status"].(string)
		if !ok {
			return common.MessageResponse{}, errors.New("status not found in RPC response")
		}

		var code int
		parsedCode, ok := response["code"].(float64)
		if !ok {
			return common.MessageResponse{}, errors.New("code not found in RPC response")
		}
		code = int(parsedCode)

		var data []interface{}
		if code != utils.SUCCESS {
			if msg, ok := response["data"].(string); ok {
				data = []interface{}{msg}
			} else {
				return common.MessageResponse{}, errors.New("data not found in RPC response")
			}
		} else if dataIface, ok := response["data"].([]interface{}); ok {
			data = dataIface
		} else {
			return common.MessageResponse{}, errors.New("data not found in RPC response")
		}

		return common.MessageResponse{
			Status: status,
			Code:   code,
			Data:   data,
		}, nil
	case <-time.After(10 * time.Second): // Timeout
		c.responseMap.Delete(corrID) // Clean up if timeout
		return common.MessageResponse{}, errors.New("timeout waiting for RPC response")
	}
}

// SendAckCommand sends a message without waiting for a response (Ack-based)
func (c *RabbitMQClient) SendAckCommand(queueName string, cmd string, data interface{}, isEvent bool) error {
	var message = map[string]interface{}{
		"pattern": map[string]string{
			"cmd": cmd,
		},
		"data": data,
	}

	if isEvent {
		message["pattern"] = cmd
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	log.Println("Sent Ack-based command:" + cmd)

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
