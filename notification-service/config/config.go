package config

import "os"

type Config struct {
	RabbitMQURL string
}

func LoadConfig() *Config {
	return &Config{
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://admin:admin@localhost:5672/"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}