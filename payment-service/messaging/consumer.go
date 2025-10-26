package messaging

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// PaymentProcessor defines the interface for processing payments
type PaymentProcessor interface {
	ProcessPayment(orderID, itemID string, quantity int, userEmail string) (float64, bool, string)
}

type Consumer struct {
	conn           *amqp.Connection
	channel        *amqp.Channel
	paymentService PaymentProcessor
	publisher      *Publisher
}

type OrderCreatedEvent struct {
	OrderID   string `json:"order_id"`
	ItemID    string `json:"item_id"`
	Quantity  int    `json:"quantity"`
	UserEmail string `json:"user_email"`
}

func NewConsumer(rabbitMQURL string, paymentService PaymentProcessor, publisher *Publisher) (*Consumer, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare orders exchange
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

	// Declare queue for order.created events
	queue, err := channel.QueueDeclare(
		"order.created.payment.queue",
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

	log.Println("Payment Service RabbitMQ consumer initialized successfully")

	return &Consumer{
		conn:           conn,
		channel:        channel,
		paymentService: paymentService,
		publisher:      publisher,
	}, nil
}

func (c *Consumer) Start() error {
	msgs, err := c.channel.Consume(
		"order.created.payment.queue",
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

			// Process the payment
			amount, success, message := c.paymentService.ProcessPayment(
				event.OrderID,
				event.ItemID,
				event.Quantity,
				event.UserEmail,
			)

			if success {
				// Publish payment.successful event
				if err := c.publisher.PublishPaymentProcessed(event.OrderID, event.ItemID, event.Quantity, event.UserEmail, amount, message); err != nil {
					log.Printf("Failed to publish payment.successful event: %v", err)
					msg.Nack(false, true) // Requeue
					continue
				}
			} else {
				// Publish payment.failed event
				if err := c.publisher.PublishPaymentFailed(event.OrderID, event.ItemID, event.Quantity, event.UserEmail, message); err != nil {
					log.Printf("Failed to publish payment.failed event: %v", err)
					msg.Nack(false, true) // Requeue
					continue
				}
			}

			// Acknowledge the message
			msg.Ack(false)
		}
	}()

	log.Println("Payment Consumer started, waiting for messages...")
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
