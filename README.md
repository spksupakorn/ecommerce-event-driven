# Event-Driven E-commerce Order Processing System

A microservices-based e-commerce order processing system built with Go, Gin framework, and RabbitMQ. Implements the **Saga Pattern** for distributed transaction management and failure handling.

## 🏗️ Architecture

```
┌─────────────┐      ┌─────────────┐      ┌──────────────┐
│   Client    │      │  RabbitMQ   │      │  PostgreSQL  │
└──────┬──────┘      └──────┬──────┘      └──────┬───────┘
       │                    │                    │
       │ POST /orders       │                    │
       │                    │                    │
┌──────▼──────────────────────────────────────────▼───────┐
│        Order Service (Producer & Consumer)              │
│  • Receives HTTP requests                               │
│  • Saves orders to DB (PENDING)                         │
│  • Publishes order.created events                       │
│  • Listens: inventory.failed, payment.failed → CANCEL   │
│  • Listens: inventory.successful → COMPLETE ✅          │
└─────────────────────────┬───────────────────────────────┘
                          │
                          │ order.created
                          │
┌─────────────────────────▼───────────────────────────────┐
│          Payment Service (Consumer/Producer)            │
│  • Subscribes to order.created                          │
│  • Processes payment (2-second simulation)              │
│  • Stores payment amount for refunds                    │
│  • Publishes payment.successful / payment.failed        │
│  • 💰 REFUND CONSUMER: Listens to inventory.failed     │
│  •    → Issues refund (compensation transaction)        │
│  •    → Publishes payment.refunded                      │
└─────────────┬───────────────────────────┬───────────────┘
              │                           │
              │ payment.successful        │ inventory.failed
              │                           │ (triggers REFUND)
              │                           │
┌─────────────▼───────────────────────────▼───────────────┐
│          Inventory Service (Consumer/Producer)          │
│  • Subscribes to payment.successful                     │
│  • Checks & reserves stock                              │
│  • Publishes inventory.successful (success)             │
│  • Publishes inventory.failed (out of stock)            │
└─────────────┬───────────────────────────┬───────────────┘
              │                           │
              │ inventory.successful      │ inventory.failed
              │ (complete order)          │
              │                           │
┌─────────────▼───────────────────────────▼───────────────┐
│          Notification Service (Consumer)                │
│  • Subscribes to inventory.successful → completion      │
│  • Subscribes to inventory.failed → out-of-stock        │
│  • Subscribes to payment.failed → payment failure       │
│  • Subscribes to payment.refunded → refund notice       │
│  • Sends email notifications                            │
└─────────────────────────────────────────────────────────┘

═══════════════════════════════════════════════════════════
              🔄 COMPENSATION TRANSACTION FLOW
═══════════════════════════════════════════════════════════
  Payment Successful ($500) → Inventory Fails
        ↓
  inventory.failed event
        ↓
  Payment Service (Refund Consumer)
        ↓
  Issues Refund ($500) 💰
        ↓
  payment.refunded event → Customer notified
═══════════════════════════════════════════════════════════
```

## ✨ Key Features

- **Complete Saga Pattern** - Full choreography-based saga with compensation transactions
- **Payment Refunds** - Automatic refund when inventory fails (complete compensating transaction)
- **Order Completion** - Orders transition to COMPLETED status on successful processing
- **Extended Event Chain** - Realistic multi-step processing: Order → Payment → Inventory → Completion
- **Payment Processing** - Simulated payment gateway with 2-second processing delay
- **Failure Handling** - Automatic order cancellation and refunds on inventory or payment failure
- **Event-Driven Architecture** - Asynchronous communication via RabbitMQ
- **Compensating Transactions** - Automatic rollback with refunds on failures
- **Microservices Design** - Independent, loosely coupled services
- **Database Per Service** - Each service manages its own data

## 🛠️ Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- Git

## 🚀 Quick Start

### 1. Clone the Repository

```bash
git clone <your-repo-url>
cd ecommerce-event-driven
```

### 2. Start All Services

```bash
docker-compose up --build
```

This will start:
- **Order Service** (Port 8080)
- **Payment Service** (Background)
- **Inventory Service** (Background)
- **Notification Service** (Background)
- **RabbitMQ** (Ports 5672, 15672)
- **PostgreSQL** (Order DB: 5432, Inventory DB: 5433)

### 3. Verify Services

Check service health:
```bash
curl http://localhost:8080/health
```

Access RabbitMQ Management UI:
```
URL: http://localhost:15672
Username: admin
Password: admin
```

## 📝 Testing the System

### Create an Order (Success Case)

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "item_id": "product-001",
    "quantity": 2,
    "user_email": "customer@example.com"
  }'
```

**Expected Response:**
```json
{
  "message": "Order accepted for processing",
  "order_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "PENDING"
}
```

**Expected Behavior:**
- Payment is processed (2-second delay)
- Inventory is reserved and deducted
- Order status updated to `COMPLETED`
- Customer receives completion email

### Test Out of Stock (Saga Pattern with Refund)

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "item_id": "product-001",
    "quantity": 999,
    "user_email": "customer@example.com"
  }'
```

**Expected Response:**
```json
{
  "message": "Order accepted for processing",
  "order_id": "550e8400-e29b-41d4-a716-446655440001",
  "status": "PENDING"
}
```

**Expected Saga Compensation:**
1. Payment Service processes payment (2-second delay) - **Payment captured**
2. Inventory Service detects insufficient stock
3. Publishes `inventory.failed` event
4. **Payment Service receives failure → issues REFUND (compensation)**
5. Publishes `payment.refunded` event
6. Order Service automatically updates status to `CANCELLED`
7. Customer receives refund notification email

### Check Order Status

```bash
curl http://localhost:8080/api/v1/orders/{order_id}
```

**Possible Statuses:**
- `PENDING` - Initial state, awaiting payment and inventory processing
- `COMPLETED` - Successfully processed through payment and inventory
- `CANCELLED` - Automatically cancelled due to payment or inventory failure

### Available Products

- `product-001` - Laptop (Stock: 100)
- `product-002` - Mouse (Stock: 500)
- `product-003` - Keyboard (Stock: 300)

## 📊 Monitoring Logs

### Watch Order Service Logs
```bash
docker-compose logs -f order-service
```

### Watch Payment Service Logs
```bash
docker-compose logs -f payment-service
```

### Watch Inventory Service Logs
```bash
docker-compose logs -f inventory-service
```

### Watch Notification Service Logs
```bash
docker-compose logs -f notification-service
```

### Watch All Services
```bash
docker-compose logs -f
```

## 🔍 Event Flow

### Success Flow (Complete Saga)
1. **Client** → POST /orders → **Order Service**
2. **Order Service** → Saves order (PENDING) → **PostgreSQL**
3. **Order Service** → Publishes `order.created` → **RabbitMQ**
4. **RabbitMQ** → Delivers event → **Payment Service**
5. **Payment Service** → Processes payment (2-second delay) → **Stores payment data**
6. **Payment Service** → Publishes `payment.successful` → **RabbitMQ**
7. **RabbitMQ** → Delivers event → **Inventory Service**
8. **Inventory Service** → Checks & reserves stock → **PostgreSQL**
9. **Inventory Service** → Publishes `inventory.successful` → **RabbitMQ**
10. **RabbitMQ** → Delivers to:
    - **Order Service** → Updates order status to **COMPLETED** ✅
    - **Notification Service** → Sends order completion email

### Failure Flow with Compensation (Saga Pattern - Out of Stock)
1. **Client** → POST /orders → **Order Service**
2. **Order Service** → Saves order (PENDING) → **PostgreSQL**
3. **Order Service** → Publishes `order.created` → **RabbitMQ**
4. **RabbitMQ** → Delivers event → **Payment Service**
5. **Payment Service** → Processes payment (2-second delay) → **Stores payment: $500**
6. **Payment Service** → Publishes `payment.successful` → **RabbitMQ**
7. **RabbitMQ** → Delivers event → **Inventory Service**
8. **Inventory Service** → Detects insufficient stock → **PostgreSQL**
9. **Inventory Service** → Publishes `inventory.failed` → **RabbitMQ**
10. **RabbitMQ** → Delivers to multiple consumers:
    - **Payment Service** → **Issues REFUND ($500)** 💰 (Compensation Transaction)
    - **Payment Service** → Publishes `payment.refunded`
    - **Order Service** → Updates order status to CANCELLED
    - **Notification Service** → Sends refund email

### Failure Flow (Saga Pattern - Payment Failed)
1. **Client** → POST /orders → **Order Service**
2. **Order Service** → Saves order (PENDING) → **PostgreSQL**
3. **Order Service** → Publishes `order.created` → **RabbitMQ**
4. **RabbitMQ** → Delivers event → **Payment Service**
5. **Payment Service** → Payment processing fails (5% failure rate) → Simulates declined card
6. **Payment Service** → Publishes `payment.failed` → **RabbitMQ**
7. **RabbitMQ** → Delivers to:
    - **Order Service** → Updates order status to CANCELLED
    - **Notification Service** → Sends payment failure email
8. **No inventory check** - order fails early (no refund needed)

## 🏗️ Project Structure

```
├── order-service/          # HTTP API & Event Producer/Consumer
│   ├── handlers/           # HTTP request handlers
│   ├── models/             # Data models
│   ├── repository/         # Database operations
│   ├── services/           # Business logic (order status updates)
│   ├── messaging/          # RabbitMQ publisher & consumer
│   └── database/           # DB initialization
│
├── payment-service/        # Payment Processing Consumer/Producer
│   ├── services/           # Business logic (payment processing simulation)
│   └── messaging/          # Consumer & Publisher (order.created → payment.successful/failed)
│
├── inventory-service/      # Stock Management Consumer/Producer
│   ├── services/           # Business logic
│   ├── repository/         # Stock operations
│   ├── messaging/          # Consumer & Publisher (payment.successful → inventory.processed/failed)
│   └── models/             # Product models
│
├── notification-service/   # Email Notification Consumer
│   ├── services/           # Email logic (success & failure notifications)
│   └── messaging/          # Consumer (multiple queues)
│
└── shared/                 # Shared types
    └── events/             # Event definitions (order.created, payment.successful, etc.)
```

## 🎯 Best Practices Implemented

✅ **Microservices Architecture** - Independent, loosely coupled services  
✅ **Event-Driven Design** - Asynchronous communication via RabbitMQ  
✅ **Saga Pattern** - Choreography-based saga for distributed transactions  
✅ **Compensating Transactions** - Automatic order cancellation on failure  
✅ **Database Per Service** - Each service has its own database  
✅ **Interface Segregation** - Prevents import cycles, promotes clean architecture  
✅ **Graceful Shutdown** - Proper signal handling  
✅ **Health Checks** - Docker health checks & HTTP endpoints  
✅ **Connection Pooling** - Optimized database connections  
✅ **Error Handling** - Comprehensive error management  
✅ **Logging** - Structured logging throughout  
✅ **Transaction Safety** - ACID compliance for stock operations  
✅ **Message Acknowledgment** - Reliable message processing (manual ACK/NACK)  
✅ **Docker Multi-stage Builds** - Optimized container images  
✅ **Configuration Management** - Environment-based config  
✅ **Message Persistence** - Durable queues and messages  
✅ **Multiple Queue Bindings** - Services subscribe to multiple event types  

## 🧪 Testing Different Scenarios

### Scenario 1: Successful Order
```bash
# Order with available stock
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"item_id": "product-002", "quantity": 5, "user_email": "test@example.com"}'
```
**Expected:**
- Payment is processed (2-second delay)
- Stock is deducted
- Order status updated to `COMPLETED`
- Completion email sent

### Scenario 2: Insufficient Stock (Saga Pattern with Automatic Refund)
```bash
# Order exceeding available stock - triggers saga compensation with refund
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"item_id": "product-001", "quantity": 1000, "user_email": "test@example.com"}'
```
**Expected:**
- Payment is processed successfully (2-second delay, e.g., $50,000 charged)
- Inventory check fails
- **Automatic refund issued ($50,000 refunded)** 💰
- Order automatically changes to `CANCELLED`
- No stock deducted
- Refund email sent

### Scenario 3: Invalid Product
```bash
# Order for non-existent product
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"item_id": "invalid-product", "quantity": 1, "user_email": "test@example.com"}'
```
**Expected:**
- Payment is processed successfully
- Order changes to `CANCELLED` at inventory stage
- Error notification sent

### Scenario 4: Payment Failure (5% Chance)
```bash
# Order may fail at payment stage - run multiple times to test
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"item_id": "product-002", "quantity": 3, "user_email": "test@example.com"}'
```
**Expected (on payment failure):**
- Payment fails after 2-second processing
- Order automatically changes to `CANCELLED`
- No inventory check performed
- Payment failure email sent

### Observing the Saga Pattern

Watch the logs to see the saga in action:

```bash
# Terminal 1: Order Service (see order cancellation)
docker-compose logs -f order-service

# Terminal 2: Payment Service (see payment processing)
docker-compose logs -f payment-service

# Terminal 3: Inventory Service (see failure detection)
docker-compose logs -f inventory-service

# Terminal 4: Notification Service (see email notifications)
docker-compose logs -f notification-service
```

**What to look for:**
1. Order Service: "Published order.created event"
2. Payment Service: "Processing payment for order" (2-second delay)
3. Payment Service: "Payment successful" or "Payment failed"
4. Payment Service: "Received inventory.failed event for refund" (if stock fails)
5. Payment Service: "Refund successful for order" + "Published payment.refunded event"
6. Inventory Service: "Received payment.successful event"
7. Inventory Service: "Successfully processed inventory" OR "Failed to reserve stock"
8. Order Service: "Order {id} completed successfully!" OR "Order {id} cancelled"
9. Notification Service: "Order Completed" OR "Refund Processed" OR "Payment Failed"

## 🛑 Stopping the System

```bash
docker-compose down
```

To remove volumes as well:
```bash
docker-compose down -v
```

## 🔧 Local Development

### Run Order Service Locally
```bash
cd order-service
export DATABASE_URL="postgres://orderuser:orderpass@localhost:5432/orders_db?sslmode=disable"
export RABBITMQ_URL="amqp://admin:admin@localhost:5672/"
export SERVER_PORT="8080"
go run main.go
```

### Run Payment Service Locally
```bash
cd payment-service
export RABBITMQ_URL="amqp://admin:admin@localhost:5672/"
go run main.go
```

### Run Inventory Service Locally
```bash
cd inventory-service
export DATABASE_URL="postgres://inventoryuser:inventorypass@localhost:5433/inventory_db?sslmode=disable"
export RABBITMQ_URL="amqp://admin:admin@localhost:5672/"
go run main.go
```

### Run Notification Service Locally
```bash
cd notification-service
export RABBITMQ_URL="amqp://admin:admin@localhost:5672/"
go run main.go
```

## 📈 Performance Considerations

- **Concurrent Processing**: Inventory and Notification services process messages concurrently
- **Connection Pooling**: Database connections are pooled and reused
- **Message Persistence**: Messages survive RabbitMQ restarts
- **Transaction Locking**: Row-level locking prevents race conditions
- **Graceful Degradation**: Services continue operating even if others fail

## 🐛 Troubleshooting

### Service won't start
```bash
docker-compose logs <service-name>
```

### Database connection issues
```bash
# Check if databases are healthy
docker-compose ps
```

### RabbitMQ issues
```bash
# Check RabbitMQ logs
docker-compose logs rabbitmq

# Access management UI
open http://localhost:15672
```

### Clear all data
```bash
docker-compose down -v
docker-compose up --build
```

## 📚 Additional Resources

- [Complete Saga Pattern Documentation](COMPLETE_SAGA_PATTERN.md) - **Detailed saga with compensation transactions**
- [Original Saga Implementation](SAGA_PATTERN_IMPLEMENTATION.md) - Initial saga implementation guide
- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
- [Gin Framework](https://gin-gonic.com/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Microservices Patterns](https://microservices.io/patterns/)
- [Saga Pattern](https://microservices.io/patterns/data/saga.html)

## 🎓 Learning Objectives

This project demonstrates:
- Building microservices with Go
- **Complete Saga pattern with compensation transactions (refunds)**
- Event-driven architecture patterns
- **Automatic payment refunds on inventory failures**
- RabbitMQ topic exchanges and routing
- Database per service pattern
- Compensating transactions
- **Order lifecycle: PENDING → COMPLETED or CANCELLED**
- Message acknowledgment strategies
- Docker containerization and orchestration

## 📄 License

MIT License - Feel free to use this project for learning and development!
