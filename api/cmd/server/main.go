package main

import (
	"database/sql"
	_ "encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/BennyEisner/test-results/internal/routes"
	_ "github.com/lib/pq"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Parse port as integer
	portInt, err := strconv.Atoi(dbPort)
	if err != nil {
		portInt = 5432 // Default if parsing fails
	}

	// Create the connection string with integer port
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbHost, portInt, dbUser, dbPassword, dbName)

	var db *sql.DB
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		db, err := sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				fmt.Printf("Connected to db")
				break
			}
		}
		if i < maxRetries-1 {
			retryDelay := time.Duration(i+1) * time.Second
			fmt.Printf("Connection failed: %v. Retrying in %s...\n", err, retryDelay)
			time.Sleep(retryDelay)
		}
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}
	defer db.Close()

	addr := ":8080"
	mux := http.NewServeMux()

	// Register routes
	routes.RegisterRoutes(mux, db)

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	fmt.Println("Application running!")
	select {}
}
