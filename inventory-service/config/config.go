package config

import "os"

type Config struct {
	DatabaseURL string
	RabbitMQURL string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL: getEnv("DATABASE_URL", "postgres://inventoryuser:inventorypass@localhost:5433/inventory_db?sslmode=disable"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}