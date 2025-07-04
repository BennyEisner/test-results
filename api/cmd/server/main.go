package main

import (
	"database/sql"
	_ "encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	_ "strconv"
	_ "time"

	"github.com/BennyEisner/test-results/internal/infrastructure"
	_ "github.com/lib/pq"
	_ "go.uber.org/automaxprocs"
)

// Config holds the application configuration
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	ServerAddr string
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	portInt, err := strconv.Atoi(dbPort)
	if err != nil {
		portInt = 5432 // Default if parsing fails
	}

	return &Config{
		DBHost:     dbHost,
		DBPort:     portInt,
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		ServerAddr: ":8080",
	}
}

// connectDB establishes a database connection
func connectDB(config *Config) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test db connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Println("Connected to database")
	return db, nil
}

// createServer creates and configures the HTTP server
func createServer(db *sql.DB) http.Handler {
	// Use the clean hexagonal router - all domains migrated to hexagonal architecture
	return infrastructure.NewRouter(db)
}

// runServer starts the HTTP server
func runServer(addr string, handler http.Handler) error {
	log.Printf("Starting server on %s", addr)
	return http.ListenAndServe(addr, handler)
}

// run initializes and runs the application
func run() error {
	config := loadConfig()

	db, err := connectDB(config)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	server := createServer(db)
	return runServer(config.ServerAddr, server)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
