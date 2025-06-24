package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BennyEisner/test-results/internal/handler"
	"github.com/BennyEisner/test-results/internal/service"
	"github.com/BennyEisner/test-results/internal/utils"
)

// RegisterRoutes registers all routes to the provided ServeMux
// RegisterRoutes registers all routes to the provided ServeMux
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	// Instantiate services and handlers
	buildService := service.NewBuildService(db)
	buildHandler := handler.NewBuildHandler(buildService)

	buildExecutionService := service.NewBuildExecutionService(db)
	buildExecutionHandler := handler.NewBuildExecutionHandler(buildExecutionService)

	projectService := service.NewProjectService(db)
	projectHandler := handler.NewProjectHandler(projectService)

	testSuiteService := service.NewTestSuiteService(db)
	testSuiteHandler := handler.NewTestSuiteHandler(testSuiteService)

	testCaseService := service.NewTestCaseService(db)
	testCaseHandler := handler.NewTestCaseHandler(testCaseService)

	failuresService := service.NewFailuresService(db)
	failuresHandler := handler.NewFailuresHandler(failuresService)

	// Instantiate JUnitImportService and JUnitImportHandler
	junitImportService := service.NewJUnitImportService(db, buildService, testSuiteService, testCaseService, buildExecutionService)
	junitImportHandler := handler.NewJUnitImportHandler(junitImportService)
	// TODO: Instantiate other services and handlers as they are created

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

	// Project-related endpoints
	mux.HandleFunc("/api/db-test", projectHandler.HandleDBTest)

	mux.HandleFunc("/api/projects", projectHandler.HandleProjects)

	// Build-related endpoints
	mux.HandleFunc("/api/builds", buildHandler.HandleBuilds)

	mux.HandleFunc("/api/builds/recent", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodGet {
			buildHandler.GetRecentBuilds(w, r)
		} else {
			utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET Method is allowed")
		}
	})

	mux.HandleFunc("/api/builds/", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers to allow frontend access - anticipating direct calls
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS") // Adjusted for potential methods
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		requestPath := r.URL.Path
		trimmedPath := strings.Trim(strings.TrimPrefix(requestPath, "/api/builds/"), "/")
		parts := strings.Split(trimmedPath, "/")

		// Logging for debugging
		fmt.Printf("Router /api/builds/: requestPath='%s', trimmedPath='%s', parts=%v, len(parts)=%d\n", requestPath, trimmedPath, parts, len(parts))

		// Expected: /api/builds/{build_id} -> parts = ["{build_id}"]
		// Expected: /api/builds/{build_id}/executions -> parts = ["{build_id}", "executions"]
		// Expected: /api/builds/{build_id}/failures -> parts = ["{build_id}", "failures"]
		if len(parts) == 2 && parts[1] == "executions" && parts[0] != "" {
			fmt.Println("Router /api/builds/: Matched /api/builds/{id}/executions")
			// Route: /api/builds/{build_id}/executions
			buildExecutionHandler.HandleBuildExecutions(w, r)
		} else if len(parts) == 2 && parts[1] == "failures" && parts[0] != "" {
			fmt.Println("Router /api/builds/: Matched /api/builds/{id}/failures")
			// Route: /api/builds/{build_id}/failures
			failuresHandler.HandleBuildFailures(w, r)
		} else if len(parts) == 1 && parts[0] != "" {
			fmt.Println("Router /api/builds/: Matched /api/builds/{id}")
			// Route: /api/builds/{build_id}
			buildHandler.HandleBuildByPath(w, r)
		} else {
			fmt.Println("Router /api/builds/: No match, responding 404")
			utils.RespondWithError(w, http.StatusNotFound, "Resource not found or path malformed under /api/builds/ prefix.")
		}
	})

	// Test Suite related endpoints
	mux.HandleFunc("/api/suites/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the path is for suite test cases
		if strings.Contains(r.URL.Path, "/cases") {
			testCaseHandler.HandleSuiteTestCases(w, r)
		} else {
			testSuiteHandler.HandleTestSuiteByPath(w, r)
		}
	})

	// Test Case related endpoints
	mux.HandleFunc("/api/cases/", testCaseHandler.HandleTestCaseByPath)

	// This pattern handles /api/projects/{id}, /api/projects/{id}/suites,
	// /api/projects/{id}/suites/{suiteID}, /api/projects/{id}/suites/{suiteID}/builds,
	// and now /api/projects/{projectID}/suites/{suiteID}/junit_imports
	mux.HandleFunc("/api/projects/", func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers to allow frontend access
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		requestPath := r.URL.Path // Renamed to avoid conflict with 'path' module if ever imported
		// parts are based on the path after "/api/projects/" and with leading/trailing slashes removed from that segment.
		// e.g., for /api/projects/1/suites/2/builds/ -> parts = ["1", "suites", "2", "builds"]
		parts := strings.Split(strings.Trim(strings.TrimPrefix(requestPath, "/api/projects/"), "/"), "/")

		// Check for /api/projects/{projectID}/suites/{suiteID}/junit_imports
		if len(parts) == 4 && parts[1] == "suites" && parts[3] == "junit_imports" {
			if r.Method == http.MethodPost {
				junitImportHandler.HandleJUnitImport(w, r)
			} else {
				utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST is allowed for JUnit imports.")
			}
		} else if len(parts) == 4 && parts[1] == "suites" && parts[3] == "builds" {
			// Route: /api/projects/{projectID}/suites/{suiteID}/builds
			// This now correctly matches the expectation of HandleTestSuiteBuilds
			buildHandler.HandleTestSuiteBuilds(w, r)
		} else if len(parts) == 3 && parts[1] == "suites" {
			// Route: /api/projects/{projectID}/suites/{suiteID}
			if r.Method == http.MethodGet {
				projectID, errP := strconv.ParseInt(parts[0], 10, 64)
				suiteID, errS := strconv.ParseInt(parts[2], 10, 64)
				if errP != nil || errS != nil {
					utils.RespondWithError(w, http.StatusBadRequest, "Invalid project or suite ID format in path.")
				} else {
					testSuiteHandler.GetProjectTestSuiteByID(w, r, projectID, suiteID)
				}
			} else {
				// For other methods like PUT, DELETE on this specific path
				utils.RespondWithError(w, http.StatusMethodNotAllowed, fmt.Sprintf("Method %s not allowed on path %s. For GET, ensure IDs are valid.", r.Method, requestPath))
			}
		} else if len(parts) == 2 && parts[1] == "suites" {
			// Route: /api/projects/{projectID}/suites
			testSuiteHandler.HandleProjectTestSuites(w, r)
		} else if len(parts) == 1 && parts[0] != "" { // parts[0] is the project id, ensure it's not empty
			// Route: /api/projects/{projectID}
			projectHandler.HandleProjectByPath(w, r)
		} else {
			// Malformed path under /api/projects/ or path was just /api/projects/ (which is handled by another route)
			utils.RespondWithError(w, http.StatusNotFound, "Resource not found or path malformed under /api/projects/ prefix.")
		}
	})
}
