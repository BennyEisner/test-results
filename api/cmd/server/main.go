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

	"github.com/BennyEisner/test-results/internal/routes"
	_ "github.com/lib/pq"
	_ "go.uber.org/automaxprocs"
)

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	portInt, err := strconv.Atoi(dbPort)
	if err != nil {
		portInt = 5432 // Default if parsing fails
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, portInt, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Test db connection
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to database")

	addr := ":8080"
	router := routes.NewRouter(db)

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
