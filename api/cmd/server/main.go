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

	_ "github.com/BennyEisner/test-results/docs"
	"github.com/BennyEisner/test-results/internal/shared/container"
	_ "github.com/lib/pq"
	_ "go.uber.org/automaxprocs"
)

// @title           Test Results API
// @version         1.0
// @description     A hexagonal architecture API for managing test results, builds, and projects.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

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
	// Use the new hexagonal architecture router
	return container.NewRouter(db)
}

// runServer starts the HTTP server
func runServer(addr string, handler http.Handler) error {
	log.Printf("Starting server on %s", addr)
	err := http.ListenAndServe(addr, handler)
	log.Printf("ListenAndServe returned: %v", err)
	return err
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
