package messaging

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

// OrderStatusUpdater defines the interface for updating order status
type OrderStatusUpdater interface {
	UpdateOrderStatus(orderID string, status string) error
}

type Consumer struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	orderService OrderStatusUpdater
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

type InventorySuccessfulEvent struct {
	OrderID   string `json:"order_id"`
	ItemID    string `json:"item_id"`
	Quantity  int    `json:"quantity"`
	UserEmail string `json:"user_email"`
	Message   string `json:"message"`
}

func NewConsumer(rabbitMQURL string, orderService OrderStatusUpdater) (*Consumer, error) {
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

	// Declare queue for inventory.failed events
	inventoryQueue, err := channel.QueueDeclare(
		"inventory.failed.order.queue",
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

	// Bind inventory queue to exchange
	err = channel.QueueBind(
		inventoryQueue.Name,
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
	paymentQueue, err := channel.QueueDeclare(
		"payment.failed.order.queue",
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

	// Bind payment queue to exchange
	err = channel.QueueBind(
		paymentQueue.Name,
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

	// Declare queue for inventory.successful events
	inventorySuccessQueue, err := channel.QueueDeclare(
		"inventory.successful.order.queue",
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

	log.Println("Order Service RabbitMQ consumer initialized successfully")

	return &Consumer{
		conn:         conn,
		channel:      channel,
		orderService: orderService,
	}, nil
}

func (c *Consumer) Start() error {
	// Consume inventory.failed events
	inventoryMsgs, err := c.channel.Consume(
		"inventory.failed.order.queue",
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

	// Consume payment.failed events
	paymentMsgs, err := c.channel.Consume(
		"payment.failed.order.queue",
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

	// Consume inventory.successful events
	inventorySuccessMsgs, err := c.channel.Consume(
		"inventory.successful.order.queue",
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

	// Handle inventory.failed events
	go func() {
		for msg := range inventoryMsgs {
			var event InventoryFailedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal inventory.failed message: %v", err)
				msg.Nack(false, false) // Don't requeue
				continue
			}

			log.Printf("Received inventory.failed event: %+v", event)

			// Update order status to CANCELLED
			if err := c.orderService.UpdateOrderStatus(event.OrderID, "CANCELLED"); err != nil {
				log.Printf("Failed to update order status: %v", err)
				msg.Nack(false, true) // Requeue on failure
				continue
			}

			log.Printf("Order %s cancelled due to inventory failure: %s", event.OrderID, event.Reason)

			// Acknowledge the message
			msg.Ack(false)
		}
	}()

	// Handle payment.failed events
	go func() {
		for msg := range paymentMsgs {
			var event PaymentFailedEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal payment.failed message: %v", err)
				msg.Nack(false, false) // Don't requeue
				continue
			}

			log.Printf("Received payment.failed event: %+v", event)

			// Update order status to CANCELLED
			if err := c.orderService.UpdateOrderStatus(event.OrderID, "CANCELLED"); err != nil {
				log.Printf("Failed to update order status: %v", err)
				msg.Nack(false, true) // Requeue on failure
				continue
			}

			log.Printf("Order %s cancelled due to payment failure: %s", event.OrderID, event.Reason)

			// Acknowledge the message
			msg.Ack(false)
		}
	}()

	// Handle inventory.successful events (mark order as COMPLETED)
	go func() {
		for msg := range inventorySuccessMsgs {
			var event InventorySuccessfulEvent
			if err := json.Unmarshal(msg.Body, &event); err != nil {
				log.Printf("Failed to unmarshal inventory.successful message: %v", err)
				msg.Nack(false, false) // Don't requeue
				continue
			}

			log.Printf("Received inventory.successful event: %+v", event)

			// Update order status to COMPLETED
			if err := c.orderService.UpdateOrderStatus(event.OrderID, "COMPLETED"); err != nil {
				log.Printf("Failed to update order status: %v", err)
				msg.Nack(false, true) // Requeue on failure
				continue
			}

			log.Printf("Order %s completed successfully!", event.OrderID)

			// Acknowledge the message
			msg.Ack(false)
		}
	}()

	log.Println("Order Service consumer started, waiting for inventory.failed, payment.failed, and inventory.successful messages...")
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
