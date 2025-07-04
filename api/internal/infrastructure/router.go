package infrastructure

import (
	"net/http"
)

// Router handles HTTP routing using hexagonal architecture components
type Router struct {
	container *Container
	mux       *http.ServeMux
}

// NewRouter creates a new router with hexagonal architecture
func NewRouter(container *Container) *Router {
	router := &Router{
		container: container,
		mux:       http.NewServeMux(),
	}

	router.setupRoutes()
	return router
}

// setupRoutes configures all the routes
func (r *Router) setupRoutes() {
	// Health check endpoints
	r.mux.HandleFunc("GET /healthz", r.handleHealthCheck)
	r.mux.HandleFunc("GET /readyz", r.handleReadyCheck)

	// Project routes
	projectHandler := r.container.GetProjectHandler()
	r.mux.HandleFunc("GET /api/projects", projectHandler.GetAllProjects)
	r.mux.HandleFunc("POST /api/projects", projectHandler.CreateProject)
	r.mux.HandleFunc("GET /api/projects/{id}", projectHandler.GetProjectByID)
	r.mux.HandleFunc("PATCH /api/projects/{id}", projectHandler.UpdateProject)
	r.mux.HandleFunc("DELETE /api/projects/{id}", projectHandler.DeleteProject)

	// TestSuite routes
	testSuiteHandler := r.container.GetTestSuiteHandler()
	r.mux.HandleFunc("GET /api/projects/{id}/test-suites", testSuiteHandler.GetTestSuitesByProjectID)
	r.mux.HandleFunc("POST /api/projects/{id}/test-suites", testSuiteHandler.CreateTestSuite)
	r.mux.HandleFunc("GET /api/test-suites/{id}", testSuiteHandler.GetTestSuiteByID)
	r.mux.HandleFunc("PATCH /api/test-suites/{id}", testSuiteHandler.UpdateTestSuite)
	r.mux.HandleFunc("DELETE /api/test-suites/{id}", testSuiteHandler.DeleteTestSuite)

	// TestCase routes
	testCaseHandler := r.container.GetTestCaseHandler()
	r.mux.HandleFunc("GET /api/test-suites/{suiteID}/test-cases", testCaseHandler.GetTestCasesBySuiteID)
	r.mux.HandleFunc("POST /api/test-suites/{suiteID}/test-cases", testCaseHandler.CreateTestCase)
	r.mux.HandleFunc("GET /api/test-cases/{id}", testCaseHandler.GetTestCaseByID)
	r.mux.HandleFunc("GET /api/test-suites/{suiteID}/test-cases", testCaseHandler.GetTestCaseByName)
	r.mux.HandleFunc("PATCH /api/test-cases/{id}", testCaseHandler.UpdateTestCase)
	r.mux.HandleFunc("DELETE /api/test-cases/{id}", testCaseHandler.DeleteTestCase)

	// Build routes
	buildHandler := r.container.GetBuildHandler()
	r.mux.HandleFunc("GET /api/projects/{id}/builds", buildHandler.GetBuildsByProjectID)
	r.mux.HandleFunc("GET /api/test-suites/{id}/builds", buildHandler.GetBuildsByTestSuiteID)
	r.mux.HandleFunc("POST /api/builds", buildHandler.CreateBuild)
	r.mux.HandleFunc("GET /api/builds/{id}", buildHandler.GetBuildByID)
	r.mux.HandleFunc("PATCH /api/builds/{id}", buildHandler.UpdateBuild)
	r.mux.HandleFunc("DELETE /api/builds/{id}", buildHandler.DeleteBuild)

	// BuildTestCaseExecution routes
	buildTestCaseExecutionHandler := r.container.GetBuildTestCaseExecutionHandler()
	r.mux.HandleFunc("GET /api/builds/{id}/executions", buildTestCaseExecutionHandler.GetExecutionsByBuildID)
	r.mux.HandleFunc("POST /api/executions", buildTestCaseExecutionHandler.CreateExecution)
	r.mux.HandleFunc("GET /api/executions/{id}", buildTestCaseExecutionHandler.GetExecutionByID)
	r.mux.HandleFunc("PATCH /api/executions/{id}", buildTestCaseExecutionHandler.UpdateExecution)
	r.mux.HandleFunc("DELETE /api/executions/{id}", buildTestCaseExecutionHandler.DeleteExecution)

	// Failure routes
	failureHandler := r.container.GetFailureHandler()
	r.mux.HandleFunc("GET /api/failures/{id}", failureHandler.GetFailureByID)
	r.mux.HandleFunc("GET /api/failures/execution/{executionID}", failureHandler.GetFailureByExecutionID)
	r.mux.HandleFunc("POST /api/failures", failureHandler.CreateFailure)
	r.mux.HandleFunc("PATCH /api/failures/{id}", failureHandler.UpdateFailure)
	r.mux.HandleFunc("DELETE /api/failures/{id}", failureHandler.DeleteFailure)

	// User routes
	userHandler := r.container.GetUserHandler()
	r.mux.HandleFunc("GET /api/users/{id}", userHandler.GetUserByID)
	r.mux.HandleFunc("GET /api/users", userHandler.GetUserByUsername)
	r.mux.HandleFunc("POST /api/users", userHandler.CreateUser)
	r.mux.HandleFunc("PATCH /api/users/{id}", userHandler.UpdateUser)
	r.mux.HandleFunc("DELETE /api/users/{id}", userHandler.DeleteUser)

	// UserConfig routes
	userConfigHandler := r.container.GetUserConfigHandler()
	r.mux.HandleFunc("GET /api/users/{userID}/configs", userConfigHandler.GetUserConfig)
	r.mux.HandleFunc("POST /api/users/{userID}/configs", userConfigHandler.CreateUserConfig)
	r.mux.HandleFunc("PATCH /api/users/{userID}/configs", userConfigHandler.UpdateUserConfig)
	r.mux.HandleFunc("DELETE /api/users/{userID}/configs", userConfigHandler.DeleteUserConfig)

	// JUnit Import route
	junitImportHandler := r.container.GetJUnitImportHandler()
	r.mux.HandleFunc("POST /api/projects/{projectID}/suites/{suiteID}/junit_imports", junitImportHandler.HandleJUnitImport)

	// Note: GetProjectByName is handled in GetAllProjects when name query param is present
	// This follows REST conventions where filtering is done via query parameters
}

// ServeHTTP implements http.Handler interface
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests
	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	r.mux.ServeHTTP(w, req)
}

// handleHealthCheck handles the health check endpoint
func (router *Router) handleHealthCheck(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("OK")); err != nil {
		// Log error but don't return it since we can't modify the response at this point
		_ = err
	}
}

// handleReadyCheck handles the readiness check endpoint
func (router *Router) handleReadyCheck(w http.ResponseWriter, req *http.Request) {
	// You could add database connectivity checks here
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Ready")); err != nil {
		// Log error but don't return it since we can't modify the response at this point
		_ = err
	}
}
