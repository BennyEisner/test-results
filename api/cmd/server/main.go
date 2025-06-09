package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "testresults"
)

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})

}

var pool *sql.DB

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

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Home Page")
	})

	// Liveness probe
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	// Readiness probe
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		// Add readiness logic here (e.g. DB or cache connectivity check)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Ready")
	})

	// Metrics endpoint placeholder
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// In production, use Prometheus instrumentation here
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "# HELP dummy_metric A dummy metric\n# TYPE dummy_metric counter\ndummy_metric 1")
	})

	// Application handler (example route)
	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello from %s at %s\n", os.Getenv("HOSTNAME"), time.Now().Format(time.RFC3339))
	})

	// Database connection + query test example
	mux.HandleFunc("/api/db-test", func(w http.ResponseWriter, r *http.Request) {
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Database error: %v", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Database connection successful. Projects count: %d", count)
	})

	// GET all Projects
	mux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		//Query database for al projects
		rows, err := db.Query("SELECT id, name FROM projects ORDER BY id")
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
			return
		}
		defer rows.Close()

		projects := []Project{}
		for rows.Next() {
			var p Project
			if err := rows.Scan(&p.ID, &p.Name); err != nil {
				respondWithError(w, http.StatusInternalServerError, "Scan error: "+err.Error())
				return
			}
			projects = append(projects, p)
		}
		respondWithJSON(w, http.StatusOK, projects)
	})

	// CREATE new Project
	mux.HandleFunc("/api/projects/create", func(w http.ResponseWriter, r *http.Request) {
		//Query database for al projects
		var p Project
		decoder := json.NewDecoder(r.Body) // Creates JSON decoder to read from request body
		// Decodes JSON request body into p
		if err := decoder.Decode(&p); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
			return
		}
		defer r.Body.Close() //Closes request body

		// Ensures project name isnt empty
		if p.Name == "" {
			respondWithError(w, http.StatusBadRequest, "Project name is required")
			return
		}

		var id int                                                                                 // Variable to hold new projects ID returned from the database
		err := db.QueryRow("INSERT INTO projects(name) VALUES($1) RETURNING id", p.Name).Scan(&id) // Puts ID into id variable
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
			return
		}
		p.ID = id                                 // Updates Project struct with database ID
		respondWithJSON(w, http.StatusCreated, p) //Send success response with complete project in the response
	})

	//DESTROY Project by ID
	mux.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		// Ensure request is a DELETE request
		if r.Method != http.MethodDelete {
			respondWithError(w, http.StatusMethodNotAllowed, "Only DELETE method is allowed")
			return
		}
		//Extract the project ID from URL
		idStr := r.URL.Path[len("/api/projects/"):]
		if idStr == "" {
			respondWithError(w, http.StatusBadRequest, "Project name is required")
			return
		}
		// Convert ID string to int
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalild Project ID format")
			return
		}
		// Execute DELETE query
		result, err := db.Exec("DELETE FROM projects WHERE id = $1", id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
			return
		}
		// Essentially check if the project existed by seeing if any rows were affected
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Error checking delete result: "+err.Error())
			return
		}

		if rowsAffected == 0 {
			respondWithError(w, http.StatusNotFound, "Project not found")
			return
		}
		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Project deleted successfully"})

	})

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
