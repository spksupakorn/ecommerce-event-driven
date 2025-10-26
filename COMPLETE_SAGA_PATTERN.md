# Complete Saga Pattern Implementation with Compensation Transactions

## 📋 Overview

This document describes the **complete choreography-based Saga pattern** implementation with full compensation transactions (refunds) in the e-commerce event-driven system.

## 🎯 What is a Saga Pattern?

A **Saga** is a sequence of local transactions where each transaction updates data within a single service. If a transaction fails, the saga executes **compensating transactions** to undo the changes made by preceding transactions.

### Our Implementation

- **Pattern Type**: Choreography-based (no central orchestrator)
- **Compensation**: Automatic payment refunds on inventory failures
- **Event-Driven**: Services communicate via RabbitMQ events
- **Eventual Consistency**: Order status eventually reaches COMPLETED or CANCELLED

---

## ✅ Success Scenario (Happy Path)

### Event Flow
```
Order Created → Payment Processed → Inventory Reserved → Order Completed
```

### Detailed Steps

| Step | Service | Action | Event Published | State |
|------|---------|--------|-----------------|-------|
| 1 | Order Service | Create order in DB | `order.created` | PENDING |
| 2 | Payment Service | Process payment (2s delay) | `payment.successful` | Payment stored |
| 3 | Inventory Service | Reserve & deduct stock | `inventory.successful` | Stock updated |
| 4 | Order Service | Update order status | - | COMPLETED ✅ |
| 5 | Notification Service | Send completion email | - | Email sent |

### Key Points
- Order transitions: `PENDING` → `COMPLETED`
- Payment amount stored for potential refund
- All transactions committed
- Customer receives success email

---

## ❌ Failure Scenario 1: Payment Failure

### Event Flow
```
Order Created → Payment Failed → Order Cancelled
```

### Detailed Steps

| Step | Service | Action | Event Published | State |
|------|---------|--------|-----------------|-------|
| 1 | Order Service | Create order in DB | `order.created` | PENDING |
| 2 | Payment Service | Payment fails (5% rate) | `payment.failed` | No payment stored |
| 3 | Order Service | Cancel order | - | CANCELLED |
| 4 | Notification Service | Send failure email | - | Email sent |

### Key Points
- **Early failure** - no compensation needed
- No money charged
- Inventory never checked
- Order transitions: `PENDING` → `CANCELLED`

---

## 💰 Failure Scenario 2: Inventory Failure (WITH REFUND)

### Event Flow
```
Order Created → Payment Processed → Inventory Failed → REFUND ISSUED → Order Cancelled
```

### Detailed Steps

| Step | Service | Action | Event Published | State |
|------|---------|--------|-----------------|-------|
| 1 | Order Service | Create order in DB | `order.created` | PENDING |
| 2 | Payment Service | Charge payment ($500) | `payment.successful` | Payment stored: $500 |
| 3 | Inventory Service | Stock insufficient | `inventory.failed` | No stock change |
| 4 | **Payment Service** | **REFUND $500** 💰 | `payment.refunded` | **Compensation!** |
| 5 | Order Service | Cancel order | - | CANCELLED |
| 6 | Notification Service | Send refund email | - | Email sent |

### Key Points
- **Compensation Transaction**: Payment refund
- Money charged then refunded
- Customer receives refund email
- Order transitions: `PENDING` → `CANCELLED`
- **This completes the Saga pattern!**

---

## 🔄 Compensation Transaction Details

### What is a Compensation Transaction?

A compensating transaction **undoes** the effects of a previously committed transaction. In our system:

- **Forward Transaction**: `payment.successful` - Money charged
- **Compensating Transaction**: `payment.refunded` - Money refunded

### Implementation

```go
// Payment Service - RefundPayment method
func (s *PaymentService) RefundPayment(orderID, itemID string, quantity int, userEmail, reason string) (float64, bool, string) {
    // Retrieve original payment amount
    amount := s.payments[orderID]
    
    // Simulate refund processing (1 second)
    time.Sleep(1 * time.Second)
    
    // Remove payment from storage
    delete(s.payments, orderID)
    
    return amount, true, "Payment refunded successfully"
}
```

### Refund Consumer

The Payment Service has a dedicated **Refund Consumer** that listens for `inventory.failed` events:

```go
// payment-service/messaging/refund_consumer.go
// Listens to: inventory.failed
// Action: Issues refund
// Publishes: payment.refunded
```

---

## 📊 Event Choreography Map

### Events and Subscribers

| Event | Publisher | Subscribers | Purpose |
|-------|-----------|-------------|---------|
| `order.created` | Order Service | Payment Service | Trigger payment processing |
| `payment.successful` | Payment Service | Inventory Service | Trigger inventory check |
| `payment.failed` | Payment Service | Order Service, Notification | Cancel order early |
| `inventory.successful` | Inventory Service | Order Service, Notification | Complete order |
| `inventory.failed` | Inventory Service | **Payment Service** (refund), Order Service, Notification | Trigger compensation |
| `payment.refunded` | Payment Service | Notification Service | Notify customer of refund |

### Critical Insight

The **Payment Service subscribes to `inventory.failed`** - this is the compensation trigger!

---

## 🏗️ Architecture Decisions

### Why Choreography over Orchestration?

**Advantages:**
- ✅ No single point of failure
- ✅ Services remain loosely coupled
- ✅ Easy to add new services
- ✅ Natural event-driven design

**Trade-offs:**
- ❌ Complex flow visualization
- ❌ Distributed transaction monitoring harder
- ❌ Need careful event design

### Why Store Payment Data?

The Payment Service stores payment amounts in memory:
```go
type PaymentService struct {
    payments map[string]float64  // orderID -> amount
    mu       sync.RWMutex
}
```

**Reason**: Need original payment amount for accurate refunds!

---

## 🧪 Testing the Saga Pattern

### Test 1: Successful Order (No Compensation)
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"item_id": "product-001", "quantity": 2, "user_email": "test@example.com"}'
```

**Expected:** Order reaches `COMPLETED` status

### Test 2: Out of Stock (WITH Refund)
```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"item_id": "product-001", "quantity": 999, "user_email": "test@example.com"}'
```

**Expected:**
1. Payment processed (~$50,000)
2. Inventory fails
3. **Refund issued** (~$50,000)
4. Order `CANCELLED`

**Verify Refund:**
```bash
# Watch payment service logs
docker-compose logs -f payment-service

# Look for:
# "Processing payment for order: {id}"
# "Payment successful: $50000"
# "Received inventory.failed event for refund"
# "Refund successful for order {id}: $50000"
# "Published payment.refunded event"
```

---

## 📈 Monitoring Compensation Transactions

### Logs to Watch

**Payment Service:**
```
Processing payment for order: abc123
Payment successful for order abc123: $50000.00
Received inventory.failed event for refund: {order_id: abc123, reason: "insufficient stock"}
Processing refund for order: abc123 (reason: Inventory reservation failed: insufficient stock)
Refund successful for order abc123: $50000.00
Published payment.refunded event for order: abc123 ($50000.00 refunded)
```

**Order Service:**
```
Order abc123 cancelled due to inventory failure: insufficient stock
```

**Notification Service:**
```
Received payment.refunded event
Subject: 💰 Refund Processed - Order #abc123
Refund Amount: $50000.00
```

---

## 🎓 Key Learnings

### 1. Compensation is NOT Rollback
- Rollback: Undo uncommitted changes
- Compensation: **Undo committed changes** with new transaction

### 2. Idempotency is Critical
All event handlers should be idempotent:
- Processing same event twice = same result
- Use unique IDs for deduplication

### 3. Eventual Consistency
- Orders don't complete instantly
- 2-second payment delay + network latency
- Status changes asynchronously

### 4. Error Handling
- Manual ACK/NACK for message reliability
- Retry logic for transient failures
- Dead letter queues for persistent failures

---

## 🚀 Production Considerations

### What's Missing for Production?

1. **Distributed Tracing** - Trace requests across services
2. **Saga State Management** - Store saga state in DB
3. **Timeout Handling** - Cancel sagas that take too long
4. **Dead Letter Queues** - Handle unprocessable messages
5. **Circuit Breakers** - Prevent cascading failures
6. **Monitoring Dashboard** - Visualize saga flows
7. **Audit Log** - Track all compensation transactions

### Recommended Additions

```go
// Add correlation IDs to events
type PaymentProcessedEvent struct {
    CorrelationID string `json:"correlation_id"`
    OrderID       string `json:"order_id"`
    // ... other fields
}

// Store saga state
type SagaState struct {
    SagaID    string
    Status    string // "STARTED", "PAYMENT_DONE", "COMPLETED", "COMPENSATING"
    Steps     []Step
    CreatedAt time.Time
}
```

---

## 📚 References

- [Microservices Patterns - Chris Richardson](https://microservices.io/patterns/data/saga.html)
- [Saga Pattern - Martin Fowler](https://martinfowler.com/articles/microservices.html)
- [Event-Driven Architecture](https://martinfowler.com/articles/201701-event-driven.html)

---

## 🎯 Summary

This implementation demonstrates a **complete Saga pattern** with:

✅ **Forward Transactions**: Order → Payment → Inventory  
✅ **Compensating Transactions**: Refund on inventory failure  
✅ **Event Choreography**: No central orchestrator  
✅ **Eventual Consistency**: Order reaches final state  
✅ **Error Handling**: Automatic compensation on failures  

The key insight: **Payment Service listens to inventory.failed to trigger refunds** - this completes the compensation transaction and makes the Saga pattern complete! 🎉
