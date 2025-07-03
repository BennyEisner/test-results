package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	APIBaseURL string
}

func LoadConfig() *Config {
	// Load .env from root dir
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found") // For local testing
	}
	apiURL := os.Getenv("API_BASE_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}

	return &Config{
		APIBaseURL: apiURL,
	}
}
