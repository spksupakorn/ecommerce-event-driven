package events

import "time"

// OrderCreatedEvent represents an order creation event
type OrderCreatedEvent struct {
	OrderID   string    `json:"order_id"`
	ItemID    string    `json:"item_id"`
	Quantity  int       `json:"quantity"`
	UserEmail string    `json:"user_email"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// InventoryProcessedEvent represents an inventory processing event
type InventoryProcessedEvent struct {
	OrderID     string    `json:"order_id"`
	ItemID      string    `json:"item_id"`
	Quantity    int       `json:"quantity"`
	UserEmail   string    `json:"user_email"`
	Status      string    `json:"status"` // "SUCCESS" or "FAILED"
	Message     string    `json:"message"`
	ProcessedAt time.Time `json:"processed_at"`
}

// InventorySuccessfulEvent represents a successful inventory reservation
type InventorySuccessfulEvent struct {
	OrderID     string    `json:"order_id"`
	ItemID      string    `json:"item_id"`
	Quantity    int       `json:"quantity"`
	UserEmail   string    `json:"user_email"`
	Message     string    `json:"message"`
	ProcessedAt time.Time `json:"processed_at"`
}

// InventoryFailedEvent represents an inventory failure event (out of stock)
type InventoryFailedEvent struct {
	OrderID   string    `json:"order_id"`
	ItemID    string    `json:"item_id"`
	Quantity  int       `json:"quantity"`
	UserEmail string    `json:"user_email"`
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
}

// PaymentProcessedEvent represents a successful payment event
type PaymentProcessedEvent struct {
	OrderID     string    `json:"order_id"`
	ItemID      string    `json:"item_id"`
	Quantity    int       `json:"quantity"`
	UserEmail   string    `json:"user_email"`
	Amount      float64   `json:"amount"`
	Status      string    `json:"status"` // "SUCCESS"
	Message     string    `json:"message"`
	ProcessedAt time.Time `json:"processed_at"`
}

// PaymentFailedEvent represents a failed payment event
type PaymentFailedEvent struct {
	OrderID   string    `json:"order_id"`
	ItemID    string    `json:"item_id"`
	Quantity  int       `json:"quantity"`
	UserEmail string    `json:"user_email"`
	Reason    string    `json:"reason"`
	FailedAt  time.Time `json:"failed_at"`
}

// PaymentRefundedEvent represents a payment refund event (compensation transaction)
type PaymentRefundedEvent struct {
	OrderID    string    `json:"order_id"`
	ItemID     string    `json:"item_id"`
	Quantity   int       `json:"quantity"`
	UserEmail  string    `json:"user_email"`
	Amount     float64   `json:"amount"`
	Reason     string    `json:"reason"`
	RefundedAt time.Time `json:"refunded_at"`
}

const (
	// Event names
	EventOrderCreated        = "order.created"
	EventInventoryProcessed  = "inventory.processed"
	EventInventorySuccessful = "inventory.successful"
	EventInventoryFailed     = "inventory.failed"
	EventPaymentProcessed    = "payment.successful"
	EventPaymentFailed       = "payment.failed"
	EventPaymentRefunded     = "payment.refunded"

	// Exchange names
	ExchangeOrders    = "orders"
	ExchangeInventory = "inventory"
	ExchangePayments  = "payments"

	// Queue names
	QueueOrderCreated             = "order.created.queue"
	QueueInventoryProcessed       = "inventory.processed.queue"
	QueueInventorySuccessful      = "inventory.successful.queue"
	QueueInventorySuccessfulOrder = "inventory.successful.order.queue"
	QueueInventoryFailed          = "inventory.failed.queue"
	QueueInventoryFailedOrder     = "inventory.failed.order.queue"
	QueueInventoryFailedPayment   = "inventory.failed.payment.queue"
	QueuePaymentProcessed         = "payment.successful.queue"
	QueuePaymentFailed            = "payment.failed.queue"
	QueuePaymentRefunded          = "payment.refunded.queue"
	QueuePaymentProcessedOrder    = "payment.successful.order.queue"

	// Routing keys
	RoutingKeyOrderCreated       = "order.created"
	RoutingKeyInventoryProcessed = "inventory.processed"
	RoutingKeyInventoryFailed    = "inventory.failed"
)
