package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppName     string
	KafkaBroker string
	PostgresDSN string
}

func Load() (*Config, error) {
	cfg := &Config{
		AppName:     getEnv("APP_NAME", "notification-service"),
		KafkaBroker: getEnv("KAFKA_BROKERS", "localhost:9092"),
		PostgresDSN: os.Getenv("POSTGRES_DSN"),
	}

	if cfg.PostgresDSN == "" {
		return nil, fmt.Errorf("POSTGRES_DSN is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
