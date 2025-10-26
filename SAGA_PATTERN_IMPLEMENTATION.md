# Saga Pattern Implementation - Failure Handling

## Overview
This implementation adds the Saga Pattern for handling inventory failures (out of stock scenarios) in the event-driven e-commerce system.

## What Was Implemented

### 1. **Shared Events** (`shared/events/events.go`)
Added new event type and constants:
- `InventoryFailedEvent` - Event structure for inventory failures
- `EventInventoryFailed` - Event name constant
- `RoutingKeyInventoryFailed` - Routing key for failed events
- Queue names for different consumers

### 2. **Inventory Service** 
#### `inventory-service/messaging/publisher.go`
- Added `PublishInventoryFailed()` method to publish inventory.failed events

#### `inventory-service/services/inventory_service.go`
- Added `publishInventoryFailedEvent()` helper method
- Modified `ProcessOrder()` to publish `inventory.failed` event when stock reservation fails

### 3. **Order Service**
#### `order-service/messaging/consumer.go` (NEW)
- Created consumer to listen for `inventory.failed` events
- Uses `OrderStatusUpdater` interface to avoid import cycles
- Automatically updates order status to "CANCELLED" when inventory fails

#### `order-service/services/order_service.go` (NEW)
- Created `OrderService` with `UpdateOrderStatus()` method
- Implements the `OrderStatusUpdater` interface

#### `order-service/main.go`
- Initialized and started the consumer
- Integrated with existing HTTP server

### 4. **Notification Service**
#### `notification-service/messaging/consumer.go`
- Added binding for `inventory.failed` queue
- Added second consumer goroutine for handling failed events
- Added `InventoryFailedEvent` struct

#### `notification-service/services/notification_service.go`
- Added `SendOutOfStockNotification()` method
- Sends "Sorry, out of stock" email when inventory fails

## Event Flow (Saga Pattern)

### Success Flow:
```
1. User creates order
   ↓
2. Order Service publishes "order.created" event
   ↓
3. Inventory Service receives event & reserves stock
   ↓
4. Inventory Service publishes "inventory.processed" (SUCCESS)
   ↓
5. Notification Service sends confirmation email
```

### Failure Flow (NEW):
```
1. User creates order
   ↓
2. Order Service publishes "order.created" event
   ↓
3. Inventory Service receives event & checks stock
   ↓
4. OUT OF STOCK detected
   ↓
5. Inventory Service publishes "inventory.failed" event
   ↓
6. Order Service receives "inventory.failed"
   ├─> Updates order status to "CANCELLED"
   │
7. Notification Service receives "inventory.failed"
   └─> Sends "out of stock" email to customer
```

## Key Design Patterns Used

### 1. **Saga Pattern**
- Compensating transactions (order cancellation)
- Event-driven choreography (no central orchestrator)

### 2. **Interface Segregation**
- `OrderStatusUpdater` interface in Order Service consumer
- Prevents import cycles
- Only exposes necessary methods

### 3. **Publisher-Subscriber**
- Multiple services subscribe to same events
- Decoupled communication

## Testing the Implementation

### Test Out of Stock Scenario:
1. Start all services
2. Create an order with quantity exceeding available stock
3. Observe the logs:
   - Inventory Service: "Failed to reserve stock"
   - Inventory Service: "Published inventory.failed event"
   - Order Service: "Order X cancelled due to: insufficient stock"
   - Notification Service: "Out of Stock - Order #X Cancelled"

### Expected Behavior:
- Order status in database: `CANCELLED`
- Customer receives email notification about cancellation
- System maintains consistency across all services

## Benefits

1. **Automatic Compensation**: Failed orders are automatically cancelled
2. **Customer Notification**: Users are immediately notified of failures
3. **Data Consistency**: Order status reflects actual inventory state
4. **Decoupled Services**: Services communicate only through events
5. **Scalability**: Each service can scale independently

## Future Enhancements

Potential improvements:
- Add retry logic with exponential backoff
- Implement dead letter queues for failed messages
- Add order history/audit trail
- Implement inventory reservation timeout
- Add payment service with refund saga
