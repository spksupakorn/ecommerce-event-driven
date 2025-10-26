package services

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type PaymentService struct {
	// Store payment amounts for potential refunds
	payments map[string]float64
	mu       sync.RWMutex
}

func NewPaymentService() *PaymentService {
	return &PaymentService{
		payments: make(map[string]float64),
	}
}

// ProcessPayment simulates payment processing with a 2-second delay
func (s *PaymentService) ProcessPayment(orderID, itemID string, quantity int, userEmail string) (float64, bool, string) {
	log.Printf("Processing payment for order: %s", orderID)

	// Simulate payment processing time
	time.Sleep(2 * time.Second)

	// Calculate mock amount (random price between 10 and 1000)
	amount := float64(quantity) * (10 + rand.Float64()*990)

	// Simulate 95% success rate for payments
	// For demonstration, you can adjust this logic
	success := rand.Float64() < 0.95

	if success {
		// Store payment amount for potential refund
		s.mu.Lock()
		s.payments[orderID] = amount
		s.mu.Unlock()

		log.Printf("Payment successful for order %s: $%.2f", orderID, amount)
		return amount, true, "Payment processed successfully"
	}

	log.Printf("Payment failed for order %s", orderID)
	return 0, false, "Payment processing failed - insufficient funds or card declined"
}

// RefundPayment simulates refunding a payment (compensation transaction)
func (s *PaymentService) RefundPayment(orderID, itemID string, quantity int, userEmail, reason string) (float64, bool, string) {
	log.Printf("Processing refund for order: %s (reason: %s)", orderID, reason)

	// Retrieve the original payment amount
	s.mu.RLock()
	amount, exists := s.payments[orderID]
	s.mu.RUnlock()

	if !exists {
		log.Printf("No payment found for order %s - cannot refund", orderID)
		return 0, false, "No payment found to refund"
	}

	// Simulate refund processing time
	time.Sleep(1 * time.Second)

	// Remove payment from storage
	s.mu.Lock()
	delete(s.payments, orderID)
	s.mu.Unlock()

	log.Printf("Refund successful for order %s: $%.2f", orderID, amount)
	return amount, true, "Payment refunded successfully"
}
