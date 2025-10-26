package messaging

import (
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type PaymentProcessedEvent struct {
	OrderID     string    `json:"order_id"`
	ItemID      string    `json:"item_id"`
	Quantity    int       `json:"quantity"`
	UserEmail   string    `json:"user_email"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	ProcessedAt time.Time `json:"processed_at"`
}

type PaymentFailedEvent struct {
	OrderID   string    `json:"order_id"`
	ItemID    string    `json:"item_id"`
	Quantity  int       `json:"quantity"`
	UserEmail string    `json:"user_email"`
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
}

type PaymentRefundedEvent struct {
	OrderID    string    `json:"order_id"`
	ItemID     string    `json:"item_id"`
	Quantity   int       `json:"quantity"`
	UserEmail  string    `json:"user_email"`
	Amount     float64   `json:"amount"`
	Reason     string    `json:"reason"`
	RefundedAt time.Time `json:"refunded_at"`
}

func NewPublisher(rabbitMQURL string) (*Publisher, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
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

	log.Println("Payment Service RabbitMQ publisher initialized successfully")

	return &Publisher{
		conn:    conn,
		channel: channel,
	}, nil
}

func (p *Publisher) PublishPaymentProcessed(orderID, itemID string, quantity int, userEmail string, amount float64, message string) error {
	event := PaymentProcessedEvent{
		OrderID:     orderID,
		ItemID:      itemID,
		Quantity:    quantity,
		UserEmail:   userEmail,
		Amount:      amount,
		Status:      "SUCCESS",
		Message:     message,
		ProcessedAt: time.Now(),
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"payments",           // exchange
		"payment.successful", // routing key
		false,                // mandatory
		false,                // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Published payment.successful event for order: %s", orderID)
	return nil
}

func (p *Publisher) PublishPaymentFailed(orderID, itemID string, quantity int, userEmail string, reason string) error {
	event := PaymentFailedEvent{
		OrderID:   orderID,
		ItemID:    itemID,
		Quantity:  quantity,
		UserEmail: userEmail,
		Reason:    reason,
		FailedAt:  time.Now(),
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"payments",       // exchange
		"payment.failed", // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Published payment.failed event for order: %s (reason: %s)", orderID, reason)
	return nil
}

func (p *Publisher) PublishPaymentRefunded(orderID, itemID string, quantity int, userEmail string, amount float64, reason string) error {
	event := PaymentRefundedEvent{
		OrderID:    orderID,
		ItemID:     itemID,
		Quantity:   quantity,
		UserEmail:  userEmail,
		Amount:     amount,
		Reason:     reason,
		RefundedAt: time.Now(),
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"payments",         // exchange
		"payment.refunded", // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Published payment.refunded event for order: %s ($%.2f refunded, reason: %s)", orderID, amount, reason)
	return nil
}

func (p *Publisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
