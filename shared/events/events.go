package events

import "time"

// OrderCreatedEvent represents an order creation event
type OrderCreatedEvent struct {
	OrderID    string    `json:"order_id"`
	ItemID     string    `json:"item_id"`
	Quantity   int       `json:"quantity"`
	UserEmail  string    `json:"user_email"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

// InventoryProcessedEvent represents an inventory processing event
type InventoryProcessedEvent struct {
	OrderID       string    `json:"order_id"`
	ItemID        string    `json:"item_id"`
	Quantity      int       `json:"quantity"`
	UserEmail     string    `json:"user_email"`
	Status        string    `json:"status"` // "SUCCESS" or "FAILED"
	Message       string    `json:"message"`
	ProcessedAt   time.Time `json:"processed_at"`
}

const (
	// Event names
	EventOrderCreated       = "order.created"
	EventInventoryProcessed = "inventory.processed"

	// Exchange names
	ExchangeOrders    = "orders"
	ExchangeInventory = "inventory"

	// Queue names
	QueueOrderCreated       = "order.created.queue"
	QueueInventoryProcessed = "inventory.processed.queue"

	// Routing keys
	RoutingKeyOrderCreated       = "order.created"
	RoutingKeyInventoryProcessed = "inventory.processed"
)