package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/BennyEisner/test-results/internal/handler"
	"github.com/BennyEisner/test-results/internal/infrastructure"
	"github.com/BennyEisner/test-results/internal/middleware"
)

// NewRouter creates a router that uses only hexagonal architecture implementations
func NewRouter(sqlDB *sql.DB) http.Handler {
	mux := http.NewServeMux()

	// Create hexagonal container
	container := infrastructure.NewContainer(sqlDB)

	// Create hexagonal handlers
	hexagonalProjectHandler := handler.NewHexagonalProjectHandler(container.GetProjectService())
	hexagonalBuildHandler := handler.NewHexagonalBuildHandler(container.GetBuildService())
	hexagonalTestSuiteHandler := handler.NewHexagonalTestSuiteHandler(container.GetTestSuiteService())
	hexagonalTestCaseHandler := handler.NewHexagonalTestCaseHandler(container.GetTestCaseService())
	hexagonalBuildExecutionHandler := handler.NewHexagonalBuildExecutionHandler(container.GetBuildTestCaseExecutionService())
	hexagonalFailuresHandler := handler.NewHexagonalFailuresHandler(container.GetFailureService())

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

	// Project-related endpoints - HEXAGONAL IMPLEMENTATION
	mux.HandleFunc("GET /api/projects", hexagonalProjectHandler.GetAllProjectsHexagonal)
	mux.HandleFunc("POST /api/projects", hexagonalProjectHandler.CreateProjectHexagonal)
	mux.HandleFunc("GET /api/projects/{id}", hexagonalProjectHandler.GetProjectByIDHexagonal)
	mux.HandleFunc("PATCH /api/projects/{id}", hexagonalProjectHandler.UpdateProjectHexagonal)
	mux.HandleFunc("DELETE /api/projects/{id}", hexagonalProjectHandler.DeleteProjectHexagonal)

	// Build-related endpoints - HEXAGONAL IMPLEMENTATION
	mux.HandleFunc("GET /api/builds", hexagonalBuildHandler.GetAllBuilds)
	mux.HandleFunc("POST /api/builds", hexagonalBuildHandler.CreateBuild)
	mux.HandleFunc("GET /api/builds/{id}", hexagonalBuildHandler.GetBuildByID)
	mux.HandleFunc("PATCH /api/builds/{id}", hexagonalBuildHandler.UpdateBuild)
	mux.HandleFunc("DELETE /api/builds/{id}", hexagonalBuildHandler.DeleteBuild)

	// Test Suite related endpoints - HEXAGONAL IMPLEMENTATION
	mux.HandleFunc("GET /api/suites/{id}", hexagonalTestSuiteHandler.GetTestSuiteByID)
	mux.HandleFunc("GET /api/suites/{id}/cases", hexagonalTestCaseHandler.GetSuiteTestCases)
	mux.HandleFunc("POST /api/suites/{id}/cases", hexagonalTestCaseHandler.CreateTestCaseForSuite)

	// Test Case related endpoints - HEXAGONAL IMPLEMENTATION
	mux.HandleFunc("GET /api/cases/{id}", hexagonalTestCaseHandler.GetTestCaseByID)
	mux.HandleFunc("GET /api/test-cases/most-failed", hexagonalTestCaseHandler.GetMostFailedTests)

	// Build Execution endpoints - HEXAGONAL IMPLEMENTATION
	mux.HandleFunc("GET /api/builds/{id}/executions", hexagonalBuildExecutionHandler.GetBuildExecutions)

	// Failure endpoints - HEXAGONAL IMPLEMENTATION
	mux.HandleFunc("GET /api/builds/{id}/failures", hexagonalFailuresHandler.GetBuildFailures)

	// Nested project/suite routes - HEXAGONAL IMPLEMENTATION
	mux.HandleFunc("GET /api/projects/{id}/suites", hexagonalTestSuiteHandler.GetTestSuitesByProjectID)
	mux.HandleFunc("POST /api/projects/{id}/suites", hexagonalTestSuiteHandler.CreateTestSuiteForProject)
	mux.HandleFunc("GET /api/projects/{projectId}/suites/{suiteId}", hexagonalTestSuiteHandler.GetProjectTestSuiteByID)
	mux.HandleFunc("GET /api/projects/{projectId}/suites/{suiteId}/builds", hexagonalBuildHandler.GetBuildsByTestSuiteID)
	mux.HandleFunc("POST /api/projects/{projectId}/suites/{suiteId}/builds", hexagonalBuildHandler.CreateBuild)

	// Search function route - HEXAGONAL IMPLEMENTATION
	mux.HandleFunc("GET /api/search", container.SearchHandler.HandleSearch)

	// User Config related endpoints - HEXAGONAL IMPLEMENTATION
	mux.HandleFunc("GET /api/users/{userID}/configs", container.UserConfigHandler.GetUserConfig)
	mux.HandleFunc("POST /api/users/{userID}/configs", container.UserConfigHandler.CreateUserConfig)
	mux.HandleFunc("PUT /api/users/{userID}/configs", container.UserConfigHandler.UpdateUserConfig)
	mux.HandleFunc("DELETE /api/users/{userID}/configs", container.UserConfigHandler.DeleteUserConfig)

	// JUnit Import endpoints - HEXAGONAL IMPLEMENTATION
	mux.HandleFunc("POST /api/projects/{projectId}/suites/{suiteId}/junit_imports", container.JUnitImportHandler.HandleJUnitImport)

	// =============================================================================
	// ADDITIONAL ROUTES TO MATCH INFRASTRUCTURE ROUTER
	// =============================================================================

	// User endpoints - INFRASTRUCTURE IMPLEMENTATION
	mux.HandleFunc("GET /api/users/{id}", container.UserHandler.GetUserByID)
	mux.HandleFunc("GET /api/users", container.UserHandler.GetUserByUsername)
	mux.HandleFunc("POST /api/users", container.UserHandler.CreateUser)
	mux.HandleFunc("PATCH /api/users/{id}", container.UserHandler.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", container.UserHandler.DeleteUser)

	// Test Suite CRUD endpoints - INFRASTRUCTURE IMPLEMENTATION
	mux.HandleFunc("GET /api/projects/{id}/test-suites", container.TestSuiteHandler.GetTestSuitesByProjectID)
	mux.HandleFunc("POST /api/projects/{id}/test-suites", container.TestSuiteHandler.CreateTestSuite)
	mux.HandleFunc("GET /api/test-suites/{id}", container.TestSuiteHandler.GetTestSuiteByID)
	mux.HandleFunc("PATCH /api/test-suites/{id}", container.TestSuiteHandler.UpdateTestSuite)
	mux.HandleFunc("DELETE /api/test-suites/{id}", container.TestSuiteHandler.DeleteTestSuite)

	// Test Case CRUD endpoints - INFRASTRUCTURE IMPLEMENTATION
	mux.HandleFunc("GET /api/test-suites/{suiteID}/test-cases", container.TestCaseHandler.GetTestCasesBySuiteID)
	mux.HandleFunc("POST /api/test-suites/{suiteID}/test-cases", container.TestCaseHandler.CreateTestCase)
	mux.HandleFunc("GET /api/test-cases/{id}", container.TestCaseHandler.GetTestCaseByID)
	mux.HandleFunc("PATCH /api/test-cases/{id}", container.TestCaseHandler.UpdateTestCase)
	mux.HandleFunc("DELETE /api/test-cases/{id}", container.TestCaseHandler.DeleteTestCase)

	// Build routes by project and test suite - INFRASTRUCTURE IMPLEMENTATION
	mux.HandleFunc("GET /api/projects/{id}/builds", container.BuildHandler.GetBuildsByProjectID)
	mux.HandleFunc("GET /api/test-suites/{id}/builds", container.BuildHandler.GetBuildsByTestSuiteID)

	// Build Execution CRUD endpoints - INFRASTRUCTURE IMPLEMENTATION
	mux.HandleFunc("POST /api/executions", container.BuildTestCaseExecutionHandler.CreateExecution)
	mux.HandleFunc("GET /api/executions/{id}", container.BuildTestCaseExecutionHandler.GetExecutionByID)
	mux.HandleFunc("PATCH /api/executions/{id}", container.BuildTestCaseExecutionHandler.UpdateExecution)
	mux.HandleFunc("DELETE /api/executions/{id}", container.BuildTestCaseExecutionHandler.DeleteExecution)

	// Failure CRUD endpoints - INFRASTRUCTURE IMPLEMENTATION
	mux.HandleFunc("GET /api/failures/{id}", container.FailureHandler.GetFailureByID)
	mux.HandleFunc("GET /api/failures/execution/{executionID}", container.FailureHandler.GetFailureByExecutionID)
	mux.HandleFunc("POST /api/failures", container.FailureHandler.CreateFailure)
	mux.HandleFunc("PATCH /api/failures/{id}", container.FailureHandler.UpdateFailure)
	mux.HandleFunc("DELETE /api/failures/{id}", container.FailureHandler.DeleteFailure)

	// Apply middleware
	var finalMux http.Handler = mux
	finalMux = middleware.Cors(finalMux)

	return finalMux
}
