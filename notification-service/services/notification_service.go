package services

import (
	"fmt"
	"log"
	"time"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendOrderConfirmation(orderID, itemID string, quantity int, userEmail, status, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	if status == "SUCCESS" {
		log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		log.Printf("📧 EMAIL NOTIFICATION")
		log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		log.Printf("To: %s", userEmail)
		log.Printf("Subject: ✅ Order Confirmed - Order #%s", orderID)
		log.Printf("")
		log.Printf("Dear Customer,")
		log.Printf("")
		log.Printf("Your order has been successfully confirmed!")
		log.Printf("")
		log.Printf("Order Details:")
		log.Printf("  • Order ID: %s", orderID)
		log.Printf("  • Item ID: %s", itemID)
		log.Printf("  • Quantity: %d", quantity)
		log.Printf("  • Status: %s", status)
		log.Printf("  • Timestamp: %s", timestamp)
		log.Printf("")
		log.Printf("Thank you for your purchase!")
		log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	} else {
		log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		log.Printf("📧 EMAIL NOTIFICATION")
		log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		log.Printf("To: %s", userEmail)
		log.Printf("Subject: ❌ Order Failed - Order #%s", orderID)
		log.Printf("")
		log.Printf("Dear Customer,")
		log.Printf("")
		log.Printf("Unfortunately, we couldn't process your order.")
		log.Printf("")
		log.Printf("Order Details:")
		log.Printf("  • Order ID: %s", orderID)
		log.Printf("  • Item ID: %s", itemID)
		log.Printf("  • Quantity: %d", quantity)
		log.Printf("  • Status: %s", status)
		log.Printf("  • Reason: %s", message)
		log.Printf("  • Timestamp: %s", timestamp)
		log.Printf("")
		log.Printf("Please contact our support team for assistance.")
		log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	}

	// Simulate email sending delay
	time.Sleep(100 * time.Millisecond)

	fmt.Println() // Add spacing for readability
}