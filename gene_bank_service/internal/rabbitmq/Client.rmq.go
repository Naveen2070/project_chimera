package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

// RPCClient handles RabbitMQ request-response messaging
type RPCClient struct {
	conn       *amqp091.Connection
	channel    *amqp091.Channel
	replyQueue amqp091.Queue
}

// NewRPCClient initializes the RabbitMQ connection
func NewRPCClient(rabbitURL string) (*RPCClient, error) {
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

	log.Println("Connected to RabbitMQ server")

	return &RPCClient{
		conn:       conn,
		channel:    ch,
		replyQueue: replyQueue,
	}, nil
}

// SendCreateCommand sends a `{ cmd: "create" }` message
func (c *RPCClient) SendCreateCommand(queueName string, cmd string, data interface{}) error {
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
		log.Printf("Failed to publish command: %v", err)
		return err
	}

	log.Println("Sent command:", string(body))
	return nil
}

// CheckQueueStatus checks if the queue is reachable
func (c *RPCClient) CheckQueueStatus(queueName string) error {
	_, err := c.channel.QueueDeclarePassive(
		queueName, // name of the queue
		false,     // durable
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	if err != nil {
		log.Printf("Queue %s is not reachable: %v", queueName, err)
		return err
	}
	return nil
}

// CheckRabbitMQStatus checks if the RabbitMQ connection is active
func (c *RPCClient) CheckRabbitMQStatus() error {
	if c.conn.IsClosed() {
		log.Println("RabbitMQ connection is closed")
		return errors.New("RabbitMQ connection is closed")
	}
	return nil
}

// Close RabbitMQ connection
func (c *RPCClient) Close() {
	c.channel.Close()
	c.conn.Close()
}
