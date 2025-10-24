package messaging

import (
	"encoding/json"

	"log"

	"github.com/streadway/amqp"
)

// OrderProcessor defines the interface for processing orders
type OrderProcessor interface {
	ProcessOrder(orderID, itemID string, quantity int, userEmail string)
}

type Consumer struct {
	conn             *amqp.Connection
	channel          *amqp.Channel
	inventoryService OrderProcessor
}

type OrderCreatedEvent struct {
	OrderID   string `json:"order_id"`
	ItemID    string `json:"item_id"`
	Quantity  int    `json:"quantity"`
	UserEmail string `json:"user_email"`
}

func NewConsumer(rabbitMQURL string, inventoryService OrderProcessor) (*Consumer, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare exchange
	err = channel.ExchangeDeclare(
		"orders",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	// Declare queue
	queue, err := channel.QueueDeclare(
		"order.created.queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	// Bind queue to exchange
	err = channel.QueueBind(
		queue.Name,
		"order.created",
		"orders",
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	log.Println("RabbitMQ consumer initialized successfully")

	return &Consumer{
		conn:             conn,
		channel:          channel,
		inventoryService: inventoryService,
	}, nil
}

func (c *Consumer) Start() error {
	msgs, err := c.channel.Consume(
		"order.created.queue",
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			var event OrderCreatedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false) // Don't requeue
				continue
			}

			log.Printf("Received order.created event: %+v", event)

			// Process the order
			c.inventoryService.ProcessOrder(event.OrderID, event.ItemID, event.Quantity, event.UserEmail)

			// Acknowledge the message
			msg.Ack(false)
		}
	}()

	log.Println("Consumer started, waiting for messages...")
	return nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
