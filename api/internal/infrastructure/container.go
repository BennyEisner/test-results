package infrastructure

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/BennyEisner/test-results/internal/application"
	"github.com/BennyEisner/test-results/internal/domain/ports"
	"github.com/BennyEisner/test-results/internal/infrastructure/database"
	httphandler "github.com/BennyEisner/test-results/internal/infrastructure/http"
	"github.com/BennyEisner/test-results/internal/middleware"
)

// Container holds all the dependencies for the application
type Container struct {
	// Database
	DB *sql.DB

	// Repositories
	ProjectRepository                ports.ProjectRepository
	BuildRepository                  ports.BuildRepository
	TestSuiteRepository              ports.TestSuiteRepository
	TestCaseRepository               ports.TestCaseRepository
	BuildExecutionRepository         ports.BuildExecutionRepository
	BuildTestCaseExecutionRepository ports.BuildTestCaseExecutionRepository
	FailureRepository                ports.FailureRepository
	UserRepository                   ports.UserRepository
	UserConfigRepository             ports.UserConfigRepository
	SearchRepository                 ports.SearchRepository

	// Services
	ProjectService                ports.ProjectService
	BuildService                  ports.BuildService
	TestSuiteService              ports.TestSuiteService
	TestCaseService               ports.TestCaseService
	BuildExecutionService         ports.BuildExecutionService
	BuildTestCaseExecutionService ports.BuildTestCaseExecutionService
	FailureService                ports.FailureService
	UserService                   ports.UserService
	UserConfigService             ports.UserConfigService
	JUnitImportService            ports.JUnitImportService
	SearchService                 ports.SearchService

	// HTTP Handlers
	ProjectHandler                *httphandler.ProjectHandler
	BuildHandler                  *httphandler.BuildHandler
	TestSuiteHandler              *httphandler.TestSuiteHandler
	TestCaseHandler               *httphandler.TestCaseHandler
	BuildExecutionHandler         *httphandler.BuildExecutionHandler
	BuildTestCaseExecutionHandler *httphandler.BuildTestCaseExecutionHandler
	FailureHandler                *httphandler.FailureHandler
	UserHandler                   *httphandler.UserHandler
	UserConfigHandler             *httphandler.UserConfigHandler
	JUnitImportHandler            *httphandler.JUnitImportHandler
	SearchHandler                 *httphandler.SearchHandler
}

// NewContainer creates a new container with all dependencies
func NewContainer(db *sql.DB) *Container {
	container := &Container{
		DB: db,
	}

	// Initialize repositories
	container.ProjectRepository = database.NewSQLProjectRepository(db)
	container.BuildRepository = database.NewSQLBuildRepository(db)
	container.TestSuiteRepository = database.NewSQLTestSuiteRepository(db)
	container.TestCaseRepository = database.NewSQLTestCaseRepository(db)
	container.BuildExecutionRepository = database.NewSQLBuildExecutionRepository(db)
	container.BuildTestCaseExecutionRepository = database.NewSQLBuildTestCaseExecutionRepository(db)
	container.FailureRepository = database.NewSQLFailureRepository(db)
	container.UserRepository = database.NewSQLUserRepository(db)
	container.UserConfigRepository = database.NewSQLUserConfigRepository(db)
	container.SearchRepository = database.NewSQLSearchRepository(db)

	// Initialize services
	container.ProjectService = application.NewProjectService(container.ProjectRepository)
	container.BuildService = application.NewBuildService(container.BuildRepository)
	container.TestSuiteService = application.NewTestSuiteService(container.TestSuiteRepository)
	container.TestCaseService = application.NewTestCaseService(container.TestCaseRepository)
	container.BuildTestCaseExecutionService = application.NewBuildTestCaseExecutionService(container.BuildTestCaseExecutionRepository)
	container.BuildExecutionService = application.NewBuildExecutionServiceAdapter(container.BuildTestCaseExecutionService)
	container.FailureService = application.NewFailureService(container.FailureRepository)
	container.UserService = application.NewUserService(container.UserRepository)
	container.UserConfigService = application.NewUserConfigService(container.UserConfigRepository)
	container.JUnitImportService = application.NewJUnitImportService(
		container.BuildService,
		container.TestSuiteService,
		container.TestCaseService,
		container.BuildExecutionService,
	)
	container.SearchService = application.NewSearchService(container.SearchRepository)

	// Initialize HTTP handlers
	container.ProjectHandler = httphandler.NewProjectHandler(container.ProjectService)
	container.BuildHandler = httphandler.NewBuildHandler(container.BuildService)
	container.TestSuiteHandler = httphandler.NewTestSuiteHandler(container.TestSuiteService)
	container.TestCaseHandler = httphandler.NewTestCaseHandler(container.TestCaseService)
	container.BuildExecutionHandler = httphandler.NewBuildExecutionHandler(container.BuildExecutionService)
	container.BuildTestCaseExecutionHandler = httphandler.NewBuildTestCaseExecutionHandler(container.BuildTestCaseExecutionService)
	container.FailureHandler = httphandler.NewFailureHandler(container.FailureService)
	container.UserHandler = httphandler.NewUserHandler(container.UserService)
	container.UserConfigHandler = httphandler.NewUserConfigHandler(container.UserConfigService)
	container.JUnitImportHandler = httphandler.NewJUnitImportHandler(container.JUnitImportService)
	container.SearchHandler = httphandler.NewSearchHandler(container.SearchService)

	return container
}

// GetProjectRepository returns the project repository
func (c *Container) GetProjectRepository() ports.ProjectRepository {
	return c.ProjectRepository
}

// GetBuildRepository returns the build repository
func (c *Container) GetBuildRepository() ports.BuildRepository {
	return c.BuildRepository
}

// GetTestSuiteRepository returns the test suite repository
func (c *Container) GetTestSuiteRepository() ports.TestSuiteRepository {
	return c.TestSuiteRepository
}

// GetTestCaseRepository returns the test case repository
func (c *Container) GetTestCaseRepository() ports.TestCaseRepository {
	return c.TestCaseRepository
}

// GetBuildExecutionRepository returns the build execution repository
func (c *Container) GetBuildExecutionRepository() ports.BuildExecutionRepository {
	return c.BuildExecutionRepository
}

// GetBuildTestCaseExecutionRepository returns the build test case execution repository
func (c *Container) GetBuildTestCaseExecutionRepository() ports.BuildTestCaseExecutionRepository {
	return c.BuildTestCaseExecutionRepository
}

// GetFailureRepository returns the failure repository
func (c *Container) GetFailureRepository() ports.FailureRepository {
	return c.FailureRepository
}

// GetUserRepository returns the user repository
func (c *Container) GetUserRepository() ports.UserRepository {
	return c.UserRepository
}

// GetUserConfigRepository returns the user config repository
func (c *Container) GetUserConfigRepository() ports.UserConfigRepository {
	return c.UserConfigRepository
}

// GetSearchRepository returns the search repository
func (c *Container) GetSearchRepository() ports.SearchRepository {
	return c.SearchRepository
}

// GetProjectService returns the project service
func (c *Container) GetProjectService() ports.ProjectService {
	return c.ProjectService
}

// GetBuildService returns the build service
func (c *Container) GetBuildService() ports.BuildService {
	return c.BuildService
}

// GetTestSuiteService returns the test suite service
func (c *Container) GetTestSuiteService() ports.TestSuiteService {
	return c.TestSuiteService
}

// GetTestCaseService returns the test case service
func (c *Container) GetTestCaseService() ports.TestCaseService {
	return c.TestCaseService
}

// GetBuildExecutionService returns the build execution service
func (c *Container) GetBuildExecutionService() ports.BuildExecutionService {
	return c.BuildExecutionService
}

// GetBuildTestCaseExecutionService returns the build test case execution service
func (c *Container) GetBuildTestCaseExecutionService() ports.BuildTestCaseExecutionService {
	return c.BuildTestCaseExecutionService
}

// GetFailureService returns the failure service
func (c *Container) GetFailureService() ports.FailureService {
	return c.FailureService
}

// GetUserService returns the user service
func (c *Container) GetUserService() ports.UserService {
	return c.UserService
}

// GetUserConfigService returns the user config service
func (c *Container) GetUserConfigService() ports.UserConfigService {
	return c.UserConfigService
}

// GetJUnitImportService returns the JUnit import service
func (c *Container) GetJUnitImportService() ports.JUnitImportService {
	return c.JUnitImportService
}

// GetSearchService returns the search service
func (c *Container) GetSearchService() ports.SearchService {
	return c.SearchService
}

// GetProjectHandler returns the project handler
func (c *Container) GetProjectHandler() *httphandler.ProjectHandler {
	return c.ProjectHandler
}

// GetBuildHandler returns the build handler
func (c *Container) GetBuildHandler() *httphandler.BuildHandler {
	return c.BuildHandler
}

// GetTestSuiteHandler returns the test suite handler
func (c *Container) GetTestSuiteHandler() *httphandler.TestSuiteHandler {
	return c.TestSuiteHandler
}

// GetTestCaseHandler returns the test case handler
func (c *Container) GetTestCaseHandler() *httphandler.TestCaseHandler {
	return c.TestCaseHandler
}

// GetBuildExecutionHandler returns the build execution handler
func (c *Container) GetBuildExecutionHandler() *httphandler.BuildExecutionHandler {
	return c.BuildExecutionHandler
}

// GetBuildTestCaseExecutionHandler returns the build test case execution handler
func (c *Container) GetBuildTestCaseExecutionHandler() *httphandler.BuildTestCaseExecutionHandler {
	return c.BuildTestCaseExecutionHandler
}

// GetFailureHandler returns the failure handler
func (c *Container) GetFailureHandler() *httphandler.FailureHandler {
	return c.FailureHandler
}

// GetUserHandler returns the user handler
func (c *Container) GetUserHandler() *httphandler.UserHandler {
	return c.UserHandler
}

// GetUserConfigHandler returns the user config handler
func (c *Container) GetUserConfigHandler() *httphandler.UserConfigHandler {
	return c.UserConfigHandler
}

// GetJUnitImportHandler returns the JUnit import handler
func (c *Container) GetJUnitImportHandler() *httphandler.JUnitImportHandler {
	return c.JUnitImportHandler
}

// GetSearchHandler returns the search handler
func (c *Container) GetSearchHandler() *httphandler.SearchHandler {
	return c.SearchHandler
}

// Close closes all resources
func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}
	return nil
}

// Validate checks if all dependencies are properly initialized
func (c *Container) Validate() error {
	if err := c.validateDB(); err != nil {
		return err
	}
	if err := c.validateRepositories(); err != nil {
		return err
	}
	return nil
}

func (c *Container) validateDB() error {
	if c.DB == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return nil
}

func (c *Container) validateRepositories() error {
	if err := c.validateProjectRepositories(); err != nil {
		return err
	}
	if err := c.validateTestRepositories(); err != nil {
		return err
	}
	if err := c.validateUserRepositories(); err != nil {
		return err
	}
	if err := c.validateSearchRepository(); err != nil {
		return err
	}
	return nil
}

func (c *Container) validateProjectRepositories() error {
	if c.ProjectRepository == nil {
		return fmt.Errorf("project repository is not initialized")
	}
	if c.BuildRepository == nil {
		return fmt.Errorf("build repository is not initialized")
	}
	return nil
}

func (c *Container) validateTestRepositories() error {
	if c.TestSuiteRepository == nil {
		return fmt.Errorf("test suite repository is not initialized")
	}
	if c.TestCaseRepository == nil {
		return fmt.Errorf("test case repository is not initialized")
	}
	if c.BuildExecutionRepository == nil {
		return fmt.Errorf("build execution repository is not initialized")
	}
	if c.BuildTestCaseExecutionRepository == nil {
		return fmt.Errorf("build test case execution repository is not initialized")
	}
	if c.FailureRepository == nil {
		return fmt.Errorf("failure repository is not initialized")
	}
	return nil
}

func (c *Container) validateUserRepositories() error {
	if c.UserRepository == nil {
		return fmt.Errorf("user repository is not initialized")
	}
	if c.UserConfigRepository == nil {
		return fmt.Errorf("user config repository is not initialized")
	}
	return nil
}

func (c *Container) validateSearchRepository() error {
	if c.SearchRepository == nil {
		return fmt.Errorf("search repository is not initialized")
	}
	return nil
}

// NewRouter creates a new HTTP router using the hexagonal architecture
func NewRouter(db *sql.DB) http.Handler {
	container := NewContainer(db)
	mux := http.NewServeMux()

	// Health and monitoring endpoints
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Ready")
	})

	// Project routes
	mux.HandleFunc("GET /api/projects", container.ProjectHandler.GetAllProjects)
	mux.HandleFunc("POST /api/projects", container.ProjectHandler.CreateProject)
	mux.HandleFunc("GET /api/projects/{id}", container.ProjectHandler.GetProjectByID)
	mux.HandleFunc("PATCH /api/projects/{id}", container.ProjectHandler.UpdateProject)
	mux.HandleFunc("DELETE /api/projects/{id}", container.ProjectHandler.DeleteProject)

	// Build routes
	mux.HandleFunc("POST /api/builds", container.BuildHandler.CreateBuild)
	mux.HandleFunc("GET /api/builds/{id}", container.BuildHandler.GetBuildByID)
	mux.HandleFunc("PUT /api/builds/{id}", container.BuildHandler.UpdateBuild)
	mux.HandleFunc("DELETE /api/builds/{id}", container.BuildHandler.DeleteBuild)

	// Test Suite routes
	mux.HandleFunc("GET /api/test-suites/{id}", container.TestSuiteHandler.GetTestSuiteByID)
	mux.HandleFunc("GET /api/projects/{projectID}/test-suites", container.TestSuiteHandler.GetTestSuitesByProjectID)
	mux.HandleFunc("POST /api/projects/{projectID}/test-suites", container.TestSuiteHandler.CreateTestSuite)
	mux.HandleFunc("PATCH /api/test-suites/{id}", container.TestSuiteHandler.UpdateTestSuite)
	mux.HandleFunc("DELETE /api/test-suites/{id}", container.TestSuiteHandler.DeleteTestSuite)

	// Test Case routes
	mux.HandleFunc("GET /api/test-cases/{id}", container.TestCaseHandler.GetTestCaseByID)
	mux.HandleFunc("GET /api/test-suites/{suiteID}/test-cases", container.TestCaseHandler.GetTestCasesBySuiteID)
	mux.HandleFunc("POST /api/test-suites/{suiteID}/test-cases", container.TestCaseHandler.CreateTestCase)
	mux.HandleFunc("PATCH /api/test-cases/{id}", container.TestCaseHandler.UpdateTestCase)
	mux.HandleFunc("DELETE /api/test-cases/{id}", container.TestCaseHandler.DeleteTestCase)

	// Build Execution routes
	mux.HandleFunc("GET /api/builds/{buildID}/executions", container.BuildExecutionHandler.GetBuildExecutions)
	mux.HandleFunc("POST /api/builds/{buildID}/executions", container.BuildExecutionHandler.CreateBuildExecutions)

	// Build Test Case Execution routes
	mux.HandleFunc("GET /api/executions/{id}", container.BuildTestCaseExecutionHandler.GetExecutionByID)
	mux.HandleFunc("POST /api/executions", container.BuildTestCaseExecutionHandler.CreateExecution)
	mux.HandleFunc("PATCH /api/executions/{id}", container.BuildTestCaseExecutionHandler.UpdateExecution)
	mux.HandleFunc("DELETE /api/executions/{id}", container.BuildTestCaseExecutionHandler.DeleteExecution)

	// Failure routes
	mux.HandleFunc("GET /api/failures/{id}", container.FailureHandler.GetFailureByID)
	mux.HandleFunc("GET /api/executions/{executionID}/failure", container.FailureHandler.GetFailureByExecutionID)
	mux.HandleFunc("POST /api/failures", container.FailureHandler.CreateFailure)
	mux.HandleFunc("PATCH /api/failures/{id}", container.FailureHandler.UpdateFailure)
	mux.HandleFunc("DELETE /api/failures/{id}", container.FailureHandler.DeleteFailure)

	// User routes
	mux.HandleFunc("GET /api/users/{id}", container.UserHandler.GetUserByID)
	mux.HandleFunc("GET /api/users", container.UserHandler.GetUserByUsername)
	mux.HandleFunc("POST /api/users", container.UserHandler.CreateUser)
	mux.HandleFunc("PATCH /api/users/{id}", container.UserHandler.UpdateUser)
	mux.HandleFunc("DELETE /api/users/{id}", container.UserHandler.DeleteUser)

	// User Config routes
	mux.HandleFunc("GET /api/users/{userID}/config", container.UserConfigHandler.GetUserConfig)
	mux.HandleFunc("POST /api/users/{userID}/config", container.UserConfigHandler.CreateUserConfig)
	mux.HandleFunc("PUT /api/users/{userID}/config", container.UserConfigHandler.UpdateUserConfig)
	mux.HandleFunc("DELETE /api/users/{userID}/config", container.UserConfigHandler.DeleteUserConfig)

	// JUnit Import routes
	mux.HandleFunc("POST /api/import/junit", container.JUnitImportHandler.ProcessJUnitData)

	// Search routes
	mux.HandleFunc("GET /api/search", container.SearchHandler.Search)

	// Apply middleware
	return middleware.Cors(mux)
}
