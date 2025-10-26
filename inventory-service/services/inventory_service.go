package services

import (
	"log"
	"time"

	"github.com/spksupakorn/ecommerce-event-driven/inventory-service/messaging"
	"github.com/spksupakorn/ecommerce-event-driven/inventory-service/repository"
)

type InventoryService struct {
	repo      *repository.InventoryRepository
	publisher *messaging.Publisher
}

func NewInventoryService(repo *repository.InventoryRepository, publisher *messaging.Publisher) *InventoryService {
	return &InventoryService{
		repo:      repo,
		publisher: publisher,
	}
}

func (s *InventoryService) ProcessOrder(orderID, itemID string, quantity int, userEmail string) {
	log.Printf("Processing order: %s for item: %s, quantity: %d", orderID, itemID, quantity)

	// Check and reserve stock
	err := s.repo.ReserveStock(itemID, quantity)
	if err != nil {
		log.Printf("Failed to reserve stock: %v", err)
		// Publish inventory.failed event for out of stock
		s.publishInventoryFailedEvent(orderID, itemID, quantity, userEmail, err.Error())
		return
	}

	// Deduct stock
	err = s.repo.DeductStock(itemID, quantity)
	if err != nil {
		log.Printf("Failed to deduct stock: %v", err)
		// Publish inventory.failed event
		s.publishInventoryFailedEvent(orderID, itemID, quantity, userEmail, err.Error())
		return
	}

	log.Printf("Successfully processed inventory for order: %s", orderID)
	s.publishInventorySuccessfulEvent(orderID, itemID, quantity, userEmail, "Stock reserved and deducted successfully")
}

func (s *InventoryService) publishInventoryEvent(orderID, itemID string, quantity int, userEmail, status, message string) {
	event := map[string]interface{}{
		"order_id":     orderID,
		"item_id":      itemID,
		"quantity":     quantity,
		"user_email":   userEmail,
		"status":       status,
		"message":      message,
		"processed_at": time.Now(),
	}

	if err := s.publisher.PublishInventoryProcessed(event); err != nil {
		log.Printf("Failed to publish inventory.processed event: %v", err)
	}
}

func (s *InventoryService) publishInventorySuccessfulEvent(orderID, itemID string, quantity int, userEmail, message string) {
	event := map[string]interface{}{
		"order_id":     orderID,
		"item_id":      itemID,
		"quantity":     quantity,
		"user_email":   userEmail,
		"message":      message,
		"processed_at": time.Now(),
	}

	if err := s.publisher.PublishInventorySuccessful(event); err != nil {
		log.Printf("Failed to publish inventory.successful event: %v", err)
	}
}

func (s *InventoryService) publishInventoryFailedEvent(orderID, itemID string, quantity int, userEmail, reason string) {
	event := map[string]interface{}{
		"order_id":   orderID,
		"item_id":    itemID,
		"quantity":   quantity,
		"user_email": userEmail,
		"reason":     reason,
		"failed_at":  time.Now(),
	}

	if err := s.publisher.PublishInventoryFailed(event); err != nil {
		log.Printf("Failed to publish inventory.failed event: %v", err)
	}
}
