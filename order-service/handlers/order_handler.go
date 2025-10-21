package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spksupakorn/ecommerce-event-driven/order-service/messaging"
	"github.com/spksupakorn/ecommerce-event-driven/order-service/models"
	"github.com/spksupakorn/ecommerce-event-driven/order-service/repository"
)

type OrderHandler struct {
	repo      *repository.OrderRepository
	publisher *messaging.Publisher
}

func NewOrderHandler(repo *repository.OrderRepository, publisher *messaging.Publisher) *OrderHandler {
	return &OrderHandler{
		repo:      repo,
		publisher: publisher,
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Save order to database
	order, err := h.repo.Create(&req)
	if err != nil {
		log.Printf("Failed to create order: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	// Publish order.created event
	event := map[string]interface{}{
		"order_id":   order.ID,
		"item_id":    order.ItemID,
		"quantity":   order.Quantity,
		"user_email": order.UserEmail,
		"status":     order.Status,
		"created_at": order.CreatedAt,
	}

	if err := h.publisher.PublishOrderCreated(event); err != nil {
		log.Printf("Failed to publish order.created event: %v", err)
		// Note: We still return success to the user as the order is saved
	}

	// Return 202 Accepted
	c.JSON(http.StatusAccepted, gin.H{
		"message":  "Order accepted for processing",
		"order_id": order.ID,
		"status":   order.Status,
	})
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")

	order, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}
