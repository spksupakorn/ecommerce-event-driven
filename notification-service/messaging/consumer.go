package messaging

import (
	"encoding/json"
	"log"

	"github.com/spksupakorn/ecommerce-event-driven/notification-service/services"
	"github.com/streadway/amqp"
)

type Consumer struct {
	conn                *amqp.Connection
	channel             *amqp.Channel
	notificationService *services.NotificationService
}

type InventoryProcessedEvent struct {
	OrderID   string `json:"order_id"`
	ItemID    string `json:"item_id"`
	Quantity  int    `json:"quantity"`
	UserEmail string `json:"user_email"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

func NewConsumer(rabbitMQURL string, notificationService *services.NotificationService) (*Consumer, error) {
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
		"inventory",
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
		"inventory.processed.queue",
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
		"inventory.processed",
		"inventory",
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	log.Println("Notification consumer initialized successfully")

	return &Consumer{
		conn:                conn,
		channel:             channel,
		notificationService: notificationService,
	}, nil
}

func (c *Consumer) Start() error {
	msgs, err := c.channel.Consume(
		"inventory.processed.queue",
		"",
		false,
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
			var event InventoryProcessedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false)
				continue
			}

			log.Printf("Received inventory.processed event: %+v", event)

			// Send notification
			c.notificationService.SendOrderConfirmation(
				event.OrderID,
				event.ItemID,
				event.Quantity,
				event.UserEmail,
				event.Status,
				event.Message,
			)

			// Acknowledge the message
			msg.Ack(false)
		}
	}()

	log.Println("Notification consumer started, waiting for messages...")
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
