package config

import (
	"os"
)

type Config struct {
	DatabaseURL  string
	RabbitMQURL  string
	ServerPort   string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://orderuser:orderpass@localhost:5432/orders_db?sslmode=disable"),
		RabbitMQURL:  getEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
		ServerPort:   getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}