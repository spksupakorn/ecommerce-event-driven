# Event-Driven E-commerce Order Processing System

A microservices-based e-commerce order processing system built with Go, Gin framework, and RabbitMQ.

## 🏗️ Architecture

```
┌─────────────┐      ┌─────────────┐      ┌──────────────┐
│   Client    │      │  RabbitMQ   │      │  PostgreSQL  │
└──────┬──────┘      └──────┬──────┘      └──────┬───────┘
       │                    │                    │
       │ POST /orders       │                    │
       │                    │                    │
┌──────▼──────────────────────────────────────────▼───────┐
│              Order Service (Producer)                   │
│  • Receives HTTP requests                               │
│  • Saves orders to DB                                   │
│  • Publishes order.created events                       │
└─────────────────────────┬───────────────────────────────┘
                          │
                          │ order.created
                          │
┌─────────────────────────▼───────────────────────────────┐
│          Inventory Service (Consumer/Producer)          │
│  • Subscribes to order.created                          │
│  • Checks & reserves stock                              │
│  • Publishes inventory.processed                        │
└─────────────────────────┬───────────────────────────────┘
                          │
                          │ inventory.processed
                          │
┌─────────────────────────▼───────────────────────────────┐
│          Notification Service (Consumer)                │
│  • Subscribes to inventory.processed                    │
│  • Sends email confirmations                            │
└─────────────────────────────────────────────────────────┘
```

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

### Check Order Status

```bash
curl http://localhost:8080/api/v1/orders/{order_id}
```

### Test Insufficient Stock

```bash
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{
    "item_id": "product-001",
    "quantity": 999,
    "user_email": "customer@example.com"
  }'
```

### Available Products

- `product-001` - Laptop (Stock: 100)
- `product-002` - Mouse (Stock: 500)
- `product-003` - Keyboard (Stock: 300)

## 📊 Monitoring Logs

### Watch Order Service Logs
```bash
docker-compose logs -f order-service
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

1. **Client** → POST /orders → **Order Service**
2. **Order Service** → Saves order (PENDING) → **PostgreSQL**
3. **Order Service** → Publishes `order.created` → **RabbitMQ**
4. **RabbitMQ** → Delivers event → **Inventory Service**
5. **Inventory Service** → Checks stock → **PostgreSQL**
6. **Inventory Service** → Publishes `inventory.processed` → **RabbitMQ**
7. **RabbitMQ** → Delivers event → **Notification Service**
8. **Notification Service** → Logs email notification

## 🏗️ Project Structure

```
├── order-service/          # HTTP API & Event Producer
│   ├── handlers/           # HTTP request handlers
│   ├── models/             # Data models
│   ├── repository/         # Database operations
│   ├── messaging/          # RabbitMQ publisher
│   └── database/           # DB initialization
│
├── inventory-service/      # Stock Management Consumer
│   ├── services/           # Business logic
│   ├── repository/         # Stock operations
│   ├── messaging/          # Consumer & Publisher
│   └── models/             # Product models
│
├── notification-service/   # Email Notification Consumer
│   ├── services/           # Email logic
│   └── messaging/          # Consumer
│
└── shared/                 # Shared types
    └── events/             # Event definitions
```

## 🎯 Best Practices Implemented

✅ **Microservices Architecture** - Independent, loosely coupled services  
✅ **Event-Driven Design** - Asynchronous communication via RabbitMQ  
✅ **Database Per Service** - Each service has its own database  
✅ **Graceful Shutdown** - Proper signal handling  
✅ **Health Checks** - Docker health checks & HTTP endpoints  
✅ **Connection Pooling** - Optimized database connections  
✅ **Error Handling** - Comprehensive error management  
✅ **Logging** - Structured logging throughout  
✅ **Transaction Safety** - ACID compliance for stock operations  
✅ **Message Acknowledgment** - Reliable message processing  
✅ **Docker Multi-stage Builds** - Optimized container images  
✅ **Configuration Management** - Environment-based config  
✅ **Message Persistence** - Durable queues and messages  

## 🧪 Testing Different Scenarios

### Scenario 1: Successful Order
```bash
# Order with available stock
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"item_id": "product-002", "quantity": 5, "user_email": "test@example.com"}'
```

### Scenario 2: Insufficient Stock
```bash
# Order exceeding available stock
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"item_id": "product-001", "quantity": 1000, "user_email": "test@example.com"}'
```

### Scenario 3: Invalid Product
```bash
# Order for non-existent product
curl -X POST http://localhost:8080/api/v1/orders \
  -H "Content-Type: application/json" \
  -d '{"item_id": "invalid-product", "quantity": 1, "user_email": "test@example.com"}'
```

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

- [RabbitMQ Documentation](https://www.rabbitmq.com/documentation.html)
- [Gin Framework](https://gin-gonic.com/docs/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Microservices Patterns](https://microservices.io/patterns/)

## 📄 License

MIT License - Feel free to use this project for learning and development!
