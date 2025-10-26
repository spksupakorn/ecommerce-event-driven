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

type InventoryFailedEvent struct {
	OrderID   string `json:"order_id"`
	ItemID    string `json:"item_id"`
	Quantity  int    `json:"quantity"`
	UserEmail string `json:"user_email"`
	Reason    string `json:"reason"`
}

type PaymentFailedEvent struct {
	OrderID   string `json:"order_id"`
	ItemID    string `json:"item_id"`
	Quantity  int    `json:"quantity"`
	UserEmail string `json:"user_email"`
	Reason    string `json:"reason"`
}

type PaymentRefundedEvent struct {
	OrderID   string  `json:"order_id"`
	ItemID    string  `json:"item_id"`
	Quantity  int     `json:"quantity"`
	UserEmail string  `json:"user_email"`
	Amount    float64 `json:"amount"`
	Reason    string  `json:"reason"`
}

type InventorySuccessfulEvent struct {
	OrderID   string `json:"order_id"`
	ItemID    string `json:"item_id"`
	Quantity  int    `json:"quantity"`
	UserEmail string `json:"user_email"`
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

	// Declare payments exchange
	err = channel.ExchangeDeclare(
		"payments",
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

	// Bind queue to exchange for inventory.processed events
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

	// Declare queue for inventory.failed events
	failedQueue, err := channel.QueueDeclare(
		"inventory.failed.notification.queue",
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

	// Bind failed queue to exchange
	err = channel.QueueBind(
		failedQueue.Name,
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

	// Declare queue for payment.failed events
	paymentFailedQueue, err := channel.QueueDeclare(
		"payment.failed.notification.queue",
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

	// Bind payment failed queue to exchange
	err = channel.QueueBind(
		paymentFailedQueue.Name,
		"payment.failed",
		"payments",
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	// Declare queue for payment.refunded events
	paymentRefundedQueue, err := channel.QueueDeclare(
		"payment.refunded.notification.queue",
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

	// Bind payment refunded queue to exchange
	err = channel.QueueBind(
		paymentRefundedQueue.Name,
		"payment.refunded",
		"payments",
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, err
	}

	// Declare queue for inventory.successful events
	inventorySuccessQueue, err := channel.QueueDeclare(
		"inventory.successful.notification.queue",
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

	// Bind inventory success queue to exchange
	err = channel.QueueBind(
		inventorySuccessQueue.Name,
		"inventory.successful",
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
	// Start consumer for inventory.processed events
	processedMsgs, err := c.channel.Consume(
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

	// Start consumer for inventory.failed events
	failedMsgs, err := c.channel.Consume(
		"inventory.failed.notification.queue",
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

	// Start consumer for payment.failed events
	paymentFailedMsgs, err := c.channel.Consume(
		"payment.failed.notification.queue",
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

	// Start consumer for payment.refunded events
	paymentRefundedMsgs, err := c.channel.Consume(
		"payment.refunded.notification.queue",
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

	// Start consumer for inventory.successful events
	inventorySuccessMsgs, err := c.channel.Consume(
		"inventory.successful.notification.queue",
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

	// Handle inventory.processed events
	go func() {
		for msg := range processedMsgs {
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

	// Handle inventory.failed events
	go func() {
		for msg := range failedMsgs {
			var event InventoryFailedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false)
				continue
			}

			log.Printf("Received inventory.failed event: %+v", event)

			// Send out of stock notification
			c.notificationService.SendOutOfStockNotification(
				event.OrderID,
				event.ItemID,
				event.Quantity,
				event.UserEmail,
				event.Reason,
			)

			// Acknowledge the message
			msg.Ack(false)
		}
	}()

	// Handle payment.failed events
	go func() {
		for msg := range paymentFailedMsgs {
			var event PaymentFailedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false)
				continue
			}

			log.Printf("Received payment.failed event: %+v", event)

			// Send payment failed notification
			c.notificationService.SendPaymentFailedNotification(
				event.OrderID,
				event.ItemID,
				event.Quantity,
				event.UserEmail,
				event.Reason,
			)

			// Acknowledge the message
			msg.Ack(false)
		}
	}()

	// Handle payment.refunded events
	go func() {
		for msg := range paymentRefundedMsgs {
			var event PaymentRefundedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false)
				continue
			}

			log.Printf("Received payment.refunded event: %+v", event)

			// Send refund notification
			c.notificationService.SendRefundNotification(
				event.OrderID,
				event.ItemID,
				event.Quantity,
				event.UserEmail,
				event.Amount,
				event.Reason,
			)

			// Acknowledge the message
			msg.Ack(false)
		}
	}()

	// Handle inventory.successful events
	go func() {
		for msg := range inventorySuccessMsgs {
			var event InventorySuccessfulEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				msg.Nack(false, false)
				continue
			}

			log.Printf("Received inventory.successful event: %+v", event)

			// Send order completion notification
			c.notificationService.SendOrderCompletionNotification(
				event.OrderID,
				event.ItemID,
				event.Quantity,
				event.UserEmail,
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
