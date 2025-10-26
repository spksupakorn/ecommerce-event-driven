package messaging

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
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

	log.Println("Inventory publisher initialized successfully")

	return &Publisher{
		conn:    conn,
		channel: channel,
	}, nil
}

func (p *Publisher) PublishInventoryProcessed(event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"inventory",
		"inventory.processed",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return err
	}

	log.Printf("Published inventory.processed event: %s", string(body))
	return nil
}

func (p *Publisher) PublishInventoryFailed(event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"inventory",
		"inventory.failed",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return err
	}

	log.Printf("Published inventory.failed event: %s", string(body))
	return nil
}

func (p *Publisher) PublishInventorySuccessful(event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"inventory",
		"inventory.successful",
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
		},
	)

	if err != nil {
		return err
	}

	log.Printf("Published inventory.successful event: %s", string(body))
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
