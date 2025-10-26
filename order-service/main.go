package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spksupakorn/ecommerce-event-driven/order-service/config"
	"github.com/spksupakorn/ecommerce-event-driven/order-service/database"
	"github.com/spksupakorn/ecommerce-event-driven/order-service/handlers"
	"github.com/spksupakorn/ecommerce-event-driven/order-service/messaging"
	"github.com/spksupakorn/ecommerce-event-driven/order-service/repository"
	"github.com/spksupakorn/ecommerce-event-driven/order-service/services"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	db, err := database.InitDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB(db)

	// Initialize RabbitMQ publisher
	publisher, err := messaging.NewPublisher(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ publisher: %v", err)
	}
	defer publisher.Close()

	// Initialize repository and handler
	orderRepo := repository.NewOrderRepository(db)
	orderHandler := handlers.NewOrderHandler(orderRepo, publisher)

	// Initialize order service
	orderService := services.NewOrderService(orderRepo)

	// Initialize and start consumer for inventory.failed events
	consumer, err := messaging.NewConsumer(cfg.RabbitMQURL, orderService)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ consumer: %v", err)
	}
	defer consumer.Close()

	if err := consumer.Start(); err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	// Setup Gin router
	router := setupRouter(orderHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Order Service starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func setupRouter(orderHandler *handlers.OrderHandler) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Order endpoints
	v1 := router.Group("/api/v1")
	{
		v1.POST("/orders", orderHandler.CreateOrder)
		v1.GET("/orders/:id", orderHandler.GetOrder)
	}

	return router
}
