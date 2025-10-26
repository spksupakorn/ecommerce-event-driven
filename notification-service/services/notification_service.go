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

func (s *NotificationService) SendOutOfStockNotification(orderID, itemID string, quantity int, userEmail, reason string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📧 EMAIL NOTIFICATION")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("To: %s", userEmail)
	log.Printf("Subject: ⚠️  Out of Stock - Order #%s Cancelled", orderID)
	log.Printf("")
	log.Printf("Dear Customer,")
	log.Printf("")
	log.Printf("We're sorry, but your order cannot be fulfilled due to insufficient stock.")
	log.Printf("")
	log.Printf("Order Details:")
	log.Printf("  • Order ID: %s", orderID)
	log.Printf("  • Item ID: %s", itemID)
	log.Printf("  • Requested Quantity: %d", quantity)
	log.Printf("  • Reason: %s", reason)
	log.Printf("  • Status: CANCELLED")
	log.Printf("  • Timestamp: %s", timestamp)
	log.Printf("")
	log.Printf("Your order has been automatically cancelled.")
	log.Printf("Please try again later or contact our support team.")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Simulate email sending delay
	time.Sleep(100 * time.Millisecond)

	fmt.Println() // Add spacing for readability
}

func (s *NotificationService) SendPaymentFailedNotification(orderID, itemID string, quantity int, userEmail, reason string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📧 EMAIL NOTIFICATION")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("To: %s", userEmail)
	log.Printf("Subject: ❌ Payment Failed - Order #%s Cancelled", orderID)
	log.Printf("")
	log.Printf("Dear Customer,")
	log.Printf("")
	log.Printf("We're sorry, but your payment could not be processed.")
	log.Printf("")
	log.Printf("Order Details:")
	log.Printf("  • Order ID: %s", orderID)
	log.Printf("  • Item ID: %s", itemID)
	log.Printf("  • Quantity: %d", quantity)
	log.Printf("  • Reason: %s", reason)
	log.Printf("  • Status: CANCELLED")
	log.Printf("  • Timestamp: %s", timestamp)
	log.Printf("")
	log.Printf("Your order has been automatically cancelled.")
	log.Printf("Please check your payment method and try again.")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Simulate email sending delay
	time.Sleep(100 * time.Millisecond)

	fmt.Println() // Add spacing for readability
}

func (s *NotificationService) SendRefundNotification(orderID, itemID string, quantity int, userEmail string, amount float64, reason string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📧 EMAIL NOTIFICATION")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("To: %s", userEmail)
	log.Printf("Subject: 💰 Refund Processed - Order #%s", orderID)
	log.Printf("")
	log.Printf("Dear Customer,")
	log.Printf("")
	log.Printf("Your payment has been refunded due to order cancellation.")
	log.Printf("")
	log.Printf("Refund Details:")
	log.Printf("  • Order ID: %s", orderID)
	log.Printf("  • Item ID: %s", itemID)
	log.Printf("  • Quantity: %d", quantity)
	log.Printf("  • Refund Amount: $%.2f", amount)
	log.Printf("  • Reason: %s", reason)
	log.Printf("  • Status: REFUNDED")
	log.Printf("  • Timestamp: %s", timestamp)
	log.Printf("")
	log.Printf("The refund will be credited to your original payment method within 5-7 business days.")
	log.Printf("We apologize for any inconvenience caused.")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Simulate email sending delay
	time.Sleep(100 * time.Millisecond)

	fmt.Println() // Add spacing for readability
}

func (s *NotificationService) SendOrderCompletionNotification(orderID, itemID string, quantity int, userEmail, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("📧 EMAIL NOTIFICATION")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	log.Printf("To: %s", userEmail)
	log.Printf("Subject: 🎉 Order Completed - Order #%s", orderID)
	log.Printf("")
	log.Printf("Dear Customer,")
	log.Printf("")
	log.Printf("Great news! Your order has been successfully completed!")
	log.Printf("")
	log.Printf("Order Details:")
	log.Printf("  • Order ID: %s", orderID)
	log.Printf("  • Item ID: %s", itemID)
	log.Printf("  • Quantity: %d", quantity)
	log.Printf("  • Status: COMPLETED")
	log.Printf("  • Message: %s", message)
	log.Printf("  • Timestamp: %s", timestamp)
	log.Printf("")
	log.Printf("Your order is being prepared for shipment.")
	log.Printf("Thank you for shopping with us!")
	log.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Simulate email sending delay
	time.Sleep(100 * time.Millisecond)

	fmt.Println() // Add spacing for readability
}
