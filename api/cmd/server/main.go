package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/BennyEisner/test-results/routes"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "testresults"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

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
	mux := http.NewServeMux()

	// Register routes
	routes.RegisterRoutes(mux, db)

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
