package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spksupakorn/ecommerce-event-driven/inventory-service/config"
	"github.com/spksupakorn/ecommerce-event-driven/inventory-service/database"
	"github.com/spksupakorn/ecommerce-event-driven/inventory-service/messaging"
	"github.com/spksupakorn/ecommerce-event-driven/inventory-service/repository"
	"github.com/spksupakorn/ecommerce-event-driven/inventory-service/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB(db)

	// Initialize repository
	inventoryRepo := repository.NewInventoryRepository(db)

	// Initialize publisher
	publisher, err := messaging.NewPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize publisher: %v", err)
	}
	defer publisher.Close()

	// Initialize service
	inventoryService := services.NewInventoryService(inventoryRepo, publisher)

	// Initialize and start consumer
	consumer, err := messaging.NewConsumer(cfg.RabbitMQURL, inventoryService)
	if err != nil {
		log.Fatalf("Failed to initialize consumer: %v", err)
	}
	defer consumer.Close()

	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	log.Println("Inventory Service started successfully")

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down Inventory Service...")
}
