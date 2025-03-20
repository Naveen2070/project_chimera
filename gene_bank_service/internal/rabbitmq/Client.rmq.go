package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"log"

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
		"cmd":  cmd,
		"data": data,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	corrID := uuid.New().String()

	// Listen for response messages
	msgs, err := c.channel.Consume(
		c.replyQueue.Name, "", true, false, false, false, nil,
	)
	if err != nil {
		return nil, err
	}

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

	log.Println("Sent RPC command:", string(body))

	// Wait for response
	for msg := range msgs {
		if msg.CorrelationId == corrID {
			var response interface{}
			err := json.Unmarshal(msg.Body, &response)
			if err != nil {
				return nil, err
			}
			return response, nil
		}
	}

	return nil, errors.New("no response received")
}

// SendAckCommand sends a message without waiting for a response (Ack-based)
func (c *RabbitMQClient) SendAckCommand(queueName string, cmd string, data interface{}) error {
	message := map[string]interface{}{
		"cmd":  cmd,
		"data": data,
	}

	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = c.channel.PublishWithContext(context.Background(),
		"", queueName, false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Printf("Failed to publish Ack command: %v", err)
		return err
	}

	log.Println("Sent Ack-based command:", string(body))
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
