package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/BennyEisner/test-results/internal/handler"
)

// RegisterRoutes registers all routes to the provided ServeMux
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	// Home Page
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Home Page")
	})

	// Health and monitoring endpoints
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Ready")
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "# HELP dummy_metric A dummy metric\n# TYPE dummy_metric counter\ndummy_metric 1")
	})

	// Application endpoints
	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello from %s at %s\n", os.Getenv("HOSTNAME"), time.Now().Format(time.RFC3339))
	})

	// Project-related endpoints - note we're using closures to inject the DB connection
	mux.HandleFunc("/api/db-test", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleDBTest(w, r, db)
	})

	mux.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleProjects(w, r, db)
	})

	// Build-related endpoints
	mux.HandleFunc("/api/builds", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleBuilds(w, r, db)
	})

	mux.HandleFunc("/api/builds/", func(w http.ResponseWriter, r *http.Request) {
		handler.HandleBuildByPath(w, r, db)
	})

	// Note: This pattern will handle /api/projects/{id}/builds
	// It needs to be specific enough not to clash with /api/projects/{id}
	mux.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is for project builds
		if strings.HasSuffix(r.URL.Path, "/builds") || strings.Contains(r.URL.Path, "/builds/") {
			handler.HandleProjectBuilds(w, r, db)
		} else {
			// Otherwise, it's a regular project path
			handler.HandleProjectByPath(w, r, db)
		}
	})
}
