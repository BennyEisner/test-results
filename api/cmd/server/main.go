package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	// GET all Projects OR CREATE new project
	mux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getProjects(w, r, db)
		} else if r.Method == http.MethodPost {
			createProject(w, r, db)
		} else {
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// Handle project-specific operations (GET, DELETE, PATCH)
	mux.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		// Special case for the create endpoint which was previously separate
		if r.URL.Path == "/api/projects/create" {
			if r.Method == http.MethodPost {
				createProject(w, r, db)
			} else {
				respondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
			}
			return
		}

		// Extract ID from the path
		pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/projects/"), "/")
		if len(pathSegments) != 1 || pathSegments[0] == "" {
			respondWithError(w, http.StatusBadRequest, "Invalid project ID in URL")
			return
		}

		idStr := pathSegments[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid project ID format")
			return
		}

		// Route based on HTTP method
		switch r.Method {
		case http.MethodGet:
			getProjectByID(w, r, id, db)
		case http.MethodPatch:
			updateProject(w, r, id, db)
		case http.MethodDelete:
			deleteProject(w, r, id, db)
		default:
			respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

// GET a single project by ID
func getProjectByID(w http.ResponseWriter, r *http.Request, id int, db *sql.DB) {
	var p Project
	err := db.QueryRow("SELECT id, name FROM projects WHERE id = $1", id).Scan(&p.ID, &p.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			respondWithError(w, http.StatusNotFound, "Project not found")
		} else {
			respondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		}
		return
	}
	respondWithJSON(w, http.StatusOK, p)
}

// GET all Projects
func getProjects(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	//Query database for all projects
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
}

// CREATE new Project
func createProject(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var p Project
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	defer r.Body.Close()

	// Ensures project name isnt empty
	if p.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Project name is required")
		return
	}

	var id int
	err := db.QueryRow("INSERT INTO projects(name) VALUES($1) RETURNING id", p.Name).Scan(&id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	p.ID = id
	respondWithJSON(w, http.StatusCreated, p)
}

// DELETE project by ID
func deleteProject(w http.ResponseWriter, r *http.Request, id int, db *sql.DB) {
	result, err := db.Exec("DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

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
}

// UPDATE project by ID
func updateProject(w http.ResponseWriter, r *http.Request, id int, db *sql.DB) {
	var updateData map[string]interface{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&updateData); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	defer r.Body.Close()

	if len(updateData) == 0 {
		respondWithError(w, http.StatusBadRequest, "No fields provided for update")
		return
	}

	// Check if project exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	if !exists {
		respondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	// SQL Update from provided parameters
	updateFields := []string{}
	values := []interface{}{}
	valueIndex := 1

	//Check which fields and then add them to the update
	if name, ok := updateData["name"].(string); ok {
		if name == "" {
			respondWithError(w, http.StatusBadRequest, "Project name cannot be empty")
			return
		}
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", valueIndex))
		values = append(values, name)
		valueIndex++
	}

	if len(updateFields) == 0 {
		respondWithError(w, http.StatusBadRequest, "No valid fields provided for update")
		return
	}

	//Build and run SQL query
	query := fmt.Sprintf("UPDATE projects SET %s WHERE id=$%d RETURNING id, name",
		strings.Join(updateFields, ", "), valueIndex)
	values = append(values, id)

	// Execute the update query
	var updatedProject Project
	err = db.QueryRow(query, values...).Scan(&updatedProject.ID, &updatedProject.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Update failed: "+err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, updatedProject)
}
