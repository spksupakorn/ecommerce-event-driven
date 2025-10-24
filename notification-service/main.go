package main

import (
	"log"

	"os"
	"os/signal"
	"syscall"

	"github.com/spksupakorn/ecommerce-event-driven/notification-service/config"
	"github.com/spksupakorn/ecommerce-event-driven/notification-service/messaging"
	"github.com/spksupakorn/ecommerce-event-driven/notification-service/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize notification service
	notificationService := services.NewNotificationService()

	// Initialize consumer
	consumer, err := messaging.NewConsumer(cfg.RabbitMQURL, notificationService)
	if err != nil {
		log.Fatalf("Failed to initialize consumer: %v", err)
	}
	defer consumer.Close()

	// Start consuming messages
	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	log.Println("Notification Service started successfully")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Notification Service...")
}
