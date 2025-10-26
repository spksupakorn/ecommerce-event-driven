package services

import (
	"log"

	"github.com/spksupakorn/ecommerce-event-driven/order-service/repository"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(repo *repository.OrderRepository) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) UpdateOrderStatus(orderID string, status string) error {
	log.Printf("Updating order %s status to %s", orderID, status)

	err := s.repo.UpdateStatus(orderID, status)
	if err != nil {
		log.Printf("Failed to update order status: %v", err)
		return err
	}

	log.Printf("Successfully updated order %s to status %s", orderID, status)
	return nil
}
