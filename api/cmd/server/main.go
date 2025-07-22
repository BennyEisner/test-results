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
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	_ "go.uber.org/automaxprocs"
)

// Config holds the application configuration
type Config struct {
	DBHost        string
	DBPort        int
	DBUser        string
	DBPassword    string
	DBName        string
	ServerAddr    string
	FrontendURL   string
	GithubKey     string
	GithubSecret  string
	SessionSecret string
}

// loadConfig loads configuration from environment variables
func loadConfig() *Config {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	frontendURL := os.Getenv("FRONTEND_URL")
	githubKey := os.Getenv("GITHUB_KEY")
	githubSecret := os.Getenv("GITHUB_SECRET")
	sessionSecret := os.Getenv("SESSION_SECRET")

	portInt, err := strconv.Atoi(dbPort)
	if err != nil {
		portInt = 5432 // Default if parsing fails
	}

	return &Config{
		DBHost:        dbHost,
		DBPort:        portInt,
		DBUser:        dbUser,
		DBPassword:    dbPassword,
		DBName:        dbName,
		ServerAddr:    ":8080",
		FrontendURL:   frontendURL,
		GithubKey:     githubKey,
		GithubSecret:  githubSecret,
		SessionSecret: sessionSecret,
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
func createServer(db *sql.DB, frontendURL string) http.Handler {
	// Use the new hexagonal architecture router
	return container.NewRouter(db, frontendURL)
}

// runServer starts the HTTP server
func runServer(addr string, handler http.Handler) error {
	log.Printf("Starting server on %s", addr)
	err := http.ListenAndServe(addr, handler)
	log.Printf("ListenAndServe returned: %v", err)
	return err
}

// initGoth initializes the Goth providers
func initGoth(config *Config) {
	callbackURL := fmt.Sprintf("http://localhost%s/auth/github/callback", config.ServerAddr)
	goth.UseProviders(
		github.New(config.GithubKey, config.GithubSecret, callbackURL),
	)

	// Initialize session store
	store := sessions.NewCookieStore([]byte(config.SessionSecret))
	store.MaxAge(86400 * 30) // 30 days
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = false // Set to true in production

	gothic.Store = store
}

// run initializes and runs the application
func run() error {
	config := loadConfig()

	initGoth(config)

	db, err := connectDB(config)
	if err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}
	defer db.Close()

	server := createServer(db, config.FrontendURL)
	return runServer(config.ServerAddr, server)
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
