package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port        string
	PostgresDSN string
	Environment string
}

func Load() (*Config, error) {
	port := getEnvOrDefault("PORT", "8080")
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		return nil, fmt.Errorf("missing required environment variable: POSTGRES_DSN")
	}

	env := getEnvOrDefault("ENVIRONMENT", "development")

	return &Config{
		Port:        port,
		PostgresDSN: dsn,
		Environment: env,
	}, nil
}

func getEnvOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

