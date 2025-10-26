package messaging

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// RefundProcessor defines the interface for processing refunds
type RefundProcessor interface {
	RefundPayment(orderID, itemID string, quantity int, userEmail, reason string) (float64, bool, string)
}

type RefundConsumer struct {
	conn           *amqp.Connection
	channel        *amqp.Channel
	paymentService RefundProcessor
	publisher      *Publisher
}

type InventoryFailedEvent struct {
	OrderID   string `json:"order_id"`
	ItemID    string `json:"item_id"`
	Quantity  int    `json:"quantity"`
	UserEmail string `json:"user_email"`
	Reason    string `json:"reason"`
}

func NewRefundConsumer(rabbitMQURL string, paymentService RefundProcessor, publisher *Publisher) (*RefundConsumer, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	// Declare inventory exchange
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

	// Declare queue for inventory.failed events
	queue, err := channel.QueueDeclare(
		"inventory.failed.payment.queue",
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
		"inventory.failed",
		"inventory",
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	log.Println("Payment Service Refund Consumer initialized successfully")

	return &RefundConsumer{
		conn:           conn,
		channel:        channel,
		paymentService: paymentService,
		publisher:      publisher,
	}, nil
}

func (c *RefundConsumer) Start() error {
	msgs, err := c.channel.Consume(
		"inventory.failed.payment.queue",
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
			var event InventoryFailedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false) // Don't requeue
				continue
			}

			log.Printf("Received inventory.failed event for refund: %+v", event)

			// Process the refund (compensation transaction)
			amount, success, message := c.paymentService.RefundPayment(
				event.OrderID,
				event.ItemID,
				event.Quantity,
				event.UserEmail,
				event.Reason,
			)

			if success {
				// Publish payment.refunded event
				refundReason := "Inventory reservation failed: " + event.Reason
				if err := c.publisher.PublishPaymentRefunded(event.OrderID, event.ItemID, event.Quantity, event.UserEmail, amount, refundReason); err != nil {
					log.Printf("Failed to publish payment.refunded event: %v", err)
					msg.Nack(false, true) // Requeue
					continue
				}
			} else {
				log.Printf("Refund processing failed: %s", message)
			}

			// Acknowledge the message
			msg.Ack(false)
		}
	}()

	log.Println("Refund Consumer started, waiting for inventory.failed messages...")
	return nil
}

func (c *RefundConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
