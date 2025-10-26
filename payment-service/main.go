package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spksupakorn/ecommerce-event-driven/payment-service/config"
	"github.com/spksupakorn/ecommerce-event-driven/payment-service/messaging"
	"github.com/spksupakorn/ecommerce-event-driven/payment-service/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize payment service
	paymentService := services.NewPaymentService()

	// Initialize publisher
	publisher, err := messaging.NewPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize publisher: %v", err)
	}
	defer publisher.Close()

	// Initialize consumer
	consumer, err := messaging.NewConsumer(cfg.RabbitMQURL, paymentService, publisher)
	if err != nil {
		log.Fatalf("Failed to initialize consumer: %v", err)
	}
	defer consumer.Close()

	// Initialize refund consumer (for compensation transactions)
	refundConsumer, err := messaging.NewRefundConsumer(cfg.RabbitMQURL, paymentService, publisher)
	if err != nil {
		log.Fatalf("Failed to initialize refund consumer: %v", err)
	}
	defer refundConsumer.Close()

	// Start consuming messages
	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	// Start refund consumer
	if err := refundConsumer.Start(); err != nil {
		log.Fatalf("Failed to start refund consumer: %v", err)
	}

	log.Println("Payment Service started successfully (with refund support)")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Payment Service...")
}
