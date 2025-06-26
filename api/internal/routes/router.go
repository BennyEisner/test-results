package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BennyEisner/test-results/internal/handler"
	"github.com/BennyEisner/test-results/internal/middleware"
	"github.com/BennyEisner/test-results/internal/service"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	projectService := service.NewProjectService(db)
	projectHandler := handler.NewProjectHandler(projectService)

	buildService := service.NewBuildService(db)
	buildHandler := handler.NewBuildHandler(buildService)

	buildExecutionService := service.NewBuildExecutionService(db)
	buildExecutionHandler := handler.NewBuildExecutionHandler(buildExecutionService)

	testSuiteService := service.NewTestSuiteService(db)
	testSuiteHandler := handler.NewTestSuiteHandler(testSuiteService)

	testCaseService := service.NewTestCaseService(db)
	testCaseHandler := handler.NewTestCaseHandler(testCaseService)

	failuresService := service.NewFailuresService(db)
	failuresHandler := handler.NewFailuresHandler(failuresService)

	junitImportService := service.NewJUnitImportService(db, buildService, testSuiteService, testCaseService, buildExecutionService)
	junitImportHandler := handler.NewJUnitImportHandler(junitImportService)

	searchService := service.NewSearchService(db)
	searchHandler := handler.NewSearchHandler(searchService)

	// Health and monitoring endpoints
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Ready")
	})

	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "# HELP dummy_metric A dummy metric\n# TYPE dummy_metric counter\ndummy_metric 1")
	})

	// Application endpoints
	mux.HandleFunc("GET /api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Hello from %s at %s\n", os.Getenv("HOSTNAME"), time.Now().Format(time.RFC3339))
	})

	// Project-related endpoints
	mux.HandleFunc("GET /api/db-test", projectHandler.HandleDBTest)
	mux.HandleFunc("GET /api/projects", projectHandler.GetProjects)
	mux.HandleFunc("POST /api/projects", projectHandler.CreateProject)
	mux.HandleFunc("GET /api/projects/{id}", projectHandler.GetProjectByID)
	mux.HandleFunc("PATCH /api/projects/{id}", projectHandler.UpdateProject)
	mux.HandleFunc("DELETE /api/projects/{id}", projectHandler.DeleteProject)

	// Build-related endpoints
	mux.HandleFunc("GET /api/builds", buildHandler.GetAllBuilds)
	mux.HandleFunc("POST /api/builds", buildHandler.CreateBuild)
	mux.HandleFunc("GET /api/builds/recent", buildHandler.GetRecentBuilds)
	mux.HandleFunc("GET /api/builds/{id}", buildHandler.GetBuildByID)
	mux.HandleFunc("PATCH /api/builds/{id}", buildHandler.UpdateBuild)
	mux.HandleFunc("DELETE /api/builds/{id}", buildHandler.DeleteBuild)
	mux.HandleFunc("GET /api/builds/{id}/executions", buildExecutionHandler.GetBuildExecutions)
	mux.HandleFunc("GET /api/builds/{id}/failures", failuresHandler.GetBuildFailures)

	// Test Suite related endpoints
	mux.HandleFunc("GET /api/suites/{id}", testSuiteHandler.GetTestSuiteByID)
	mux.HandleFunc("GET /api/suites/{id}/cases", testCaseHandler.GetSuiteTestCases)
	mux.HandleFunc("POST /api/suites/{id}/cases", testCaseHandler.CreateTestCaseForSuite)

	// Test Case related endpoints
	mux.HandleFunc("GET /api/cases/{id}", testCaseHandler.GetTestCaseByID)

	// Nested project/suite routes
	mux.HandleFunc("GET /api/projects/{id}/suites", testSuiteHandler.GetTestSuitesByProjectID)
	mux.HandleFunc("POST /api/projects/{id}/suites", testSuiteHandler.CreateTestSuiteForProject)
	mux.HandleFunc("GET /api/projects/{projectId}/suites/{suiteId}", testSuiteHandler.GetProjectTestSuiteByID)
	mux.HandleFunc("GET /api/projects/{projectId}/suites/{suiteId}/builds", buildHandler.GetBuildsByTestSuiteID)
	mux.HandleFunc("POST /api/projects/{projectId}/suites/{suiteId}/builds", buildHandler.CreateBuildForTestSuite)
	mux.HandleFunc("POST /api/projects/{projectId}/suites/{suiteId}/junit_imports", junitImportHandler.HandleJUnitImport)

	// Search function route
	mux.HandleFunc("GET /api/search", searchHandler.HandleSearch)

	// Apply middleware
	var finalMux http.Handler = mux
	finalMux = middleware.Cors(finalMux)

	return finalMux
}
