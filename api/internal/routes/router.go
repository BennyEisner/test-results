package routes

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	dbrepo "github.com/BennyEisner/test-results/internal/db"
	"github.com/BennyEisner/test-results/internal/handler"
	"github.com/BennyEisner/test-results/internal/middleware"
	"github.com/BennyEisner/test-results/internal/service"
)

func NewRouter(sqlDB *sql.DB) http.Handler {
	mux := http.NewServeMux()

	projectRepo := dbrepo.NewSQLProjectRepository(sqlDB)
	var projectService service.ProjectServiceInterface = service.NewProjectService(projectRepo)
	projectHandler := handler.NewProjectHandler(projectService)

	buildService := service.NewBuildService(sqlDB)
	buildHandler := handler.NewBuildHandler(buildService)

	buildExecutionService := service.NewBuildExecutionService(sqlDB)
	buildExecutionHandler := handler.NewBuildExecutionHandler(buildExecutionService)

	testSuiteService := service.NewTestSuiteService(sqlDB)
	testSuiteHandler := handler.NewTestSuiteHandler(testSuiteService)

	testCaseService := service.NewTestCaseService(sqlDB)
	testCaseHandler := handler.NewTestCaseHandler(testCaseService)

	failuresService := service.NewFailuresService(sqlDB)
	failuresHandler := handler.NewFailuresHandler(failuresService)

	junitImportService := service.NewJUnitImportService(sqlDB, buildService, testSuiteService, testCaseService, buildExecutionService)
	junitImportHandler := handler.NewJUnitImportHandler(junitImportService)

	searchService := service.NewSearchService(sqlDB)
	searchHandler := handler.NewSearchHandler(searchService)

	userConfigService := service.NewUserConfigService(sqlDB)
	userConfigHandler := handler.NewUserConfigHandler(userConfigService)

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
	mux.HandleFunc("GET /api/projects/{id}/builds/recent", buildHandler.GetRecentBuilds)
	mux.HandleFunc("GET /api/builds/{id}", buildHandler.GetBuildByID)
	mux.HandleFunc("PATCH /api/builds/{id}", buildHandler.UpdateBuild)
	mux.HandleFunc("DELETE /api/builds/{id}", buildHandler.DeleteBuild)
	mux.HandleFunc("GET /api/builds/duration-trends", buildHandler.GetBuildDurationTrends)
	mux.HandleFunc("GET /api/builds/{id}/executions", buildExecutionHandler.GetBuildExecutions)
	mux.HandleFunc("GET /api/builds/{id}/failures", failuresHandler.GetBuildFailures)

	// Test Suite related endpoints
	mux.HandleFunc("GET /api/suites/{id}", testSuiteHandler.GetTestSuiteByID)
	mux.HandleFunc("GET /api/suites/{id}/cases", testCaseHandler.GetSuiteTestCases)
	mux.HandleFunc("POST /api/suites/{id}/cases", testCaseHandler.CreateTestCaseForSuite)

	// Test Case related endpoints
	mux.HandleFunc("GET /api/cases/{id}", testCaseHandler.GetTestCaseByID)
	mux.HandleFunc("GET /api/test-cases/most-failed", testCaseHandler.GetMostFailedTests)

	// Nested project/suite routes
	mux.HandleFunc("GET /api/projects/{id}/suites", testSuiteHandler.GetTestSuitesByProjectID)
	mux.HandleFunc("POST /api/projects/{id}/suites", testSuiteHandler.CreateTestSuiteForProject)
	mux.HandleFunc("GET /api/projects/{projectId}/suites/{suiteId}", testSuiteHandler.GetProjectTestSuiteByID)
	mux.HandleFunc("GET /api/projects/{projectId}/suites/{suiteId}/builds", buildHandler.GetBuildsByTestSuiteID)
	mux.HandleFunc("POST /api/projects/{projectId}/suites/{suiteId}/builds", buildHandler.CreateBuildForTestSuite)
	mux.HandleFunc("POST /api/projects/{projectId}/suites/{suiteId}/junit_imports", junitImportHandler.HandleJUnitImport)

	// Search function route
	mux.HandleFunc("GET /api/search", searchHandler.HandleSearch)

	// User Config related endpoints
	mux.HandleFunc("GET /api/users/{userId}/configs", userConfigHandler.GetUserConfig)
	mux.HandleFunc("POST /api/users/{userId}/configs", userConfigHandler.SaveUserConfig)

	// Apply middleware
	var finalMux http.Handler = mux
	finalMux = middleware.Cors(finalMux)

	return finalMux
}
