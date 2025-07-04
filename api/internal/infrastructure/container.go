package infrastructure

import (
	"database/sql"

	"github.com/BennyEisner/test-results/internal/application"
	"github.com/BennyEisner/test-results/internal/domain"
	"github.com/BennyEisner/test-results/internal/infrastructure/database"
	"github.com/BennyEisner/test-results/internal/infrastructure/http"
)

// Container holds all the dependencies for the application
type Container struct {
	// Database connection
	DB *sql.DB

	// Repositories (Secondary/Driven Adapters)
	ProjectRepository                domain.ProjectRepository
	TestSuiteRepository              domain.TestSuiteRepository
	TestCaseRepository               domain.TestCaseRepository
	BuildRepository                  domain.BuildRepository
	BuildTestCaseExecutionRepository domain.BuildTestCaseExecutionRepository
	FailureRepository                domain.FailureRepository
	UserRepository                   domain.UserRepository
	UserConfigRepository             domain.UserConfigRepository
	SearchRepository                 domain.SearchRepository

	// Application Services (Use Cases)
	ProjectService                domain.ProjectService
	TestSuiteService              domain.TestSuiteService
	TestCaseService               domain.TestCaseService
	BuildService                  domain.BuildService
	BuildTestCaseExecutionService domain.BuildTestCaseExecutionService
	FailureService                domain.FailureService
	UserService                   domain.UserService
	UserConfigService             domain.UserConfigService
	SearchService                 domain.SearchService

	// HTTP Handlers (Primary/Driving Adapters)
	ProjectHandler                *http.ProjectHandler
	TestSuiteHandler              *http.TestSuiteHandler
	TestCaseHandler               *http.TestCaseHandler
	BuildHandler                  *http.BuildHandler
	BuildTestCaseExecutionHandler *http.BuildTestCaseExecutionHandler
	FailureHandler                *http.FailureHandler
	UserHandler                   *http.UserHandler
	UserConfigHandler             *http.UserConfigHandler
	SearchHandler                 *http.SearchHandler
	JUnitImportService            domain.JUnitImportService
	JUnitImportHandler            *http.JUnitImportHandler
}

// NewContainer creates a new dependency injection container
func NewContainer(db *sql.DB) *Container {
	container := &Container{
		DB: db,
	}

	// Initialize repositories
	container.ProjectRepository = database.NewSQLProjectRepository(db)
	container.TestSuiteRepository = database.NewSQLTestSuiteRepository(db)
	container.TestCaseRepository = database.NewSQLTestCaseRepository(db)
	container.BuildRepository = database.NewSQLBuildRepository(db)
	container.BuildTestCaseExecutionRepository = database.NewSQLBuildTestCaseExecutionRepository(db)
	container.FailureRepository = database.NewSQLFailureRepository(db)
	container.UserRepository = database.NewSQLUserRepository(db)
	container.UserConfigRepository = database.NewSQLUserConfigRepository(db)
	container.SearchRepository = database.NewSQLSearchRepository(db)

	// Initialize application services
	container.ProjectService = application.NewProjectService(container.ProjectRepository)
	container.TestSuiteService = application.NewTestSuiteService(container.TestSuiteRepository)
	container.TestCaseService = application.NewTestCaseService(container.TestCaseRepository)
	container.BuildService = application.NewBuildService(container.BuildRepository)
	container.BuildTestCaseExecutionService = application.NewBuildTestCaseExecutionService(container.BuildTestCaseExecutionRepository)
	container.FailureService = application.NewFailureService(container.FailureRepository)
	container.UserService = application.NewUserService(container.UserRepository)
	container.UserConfigService = application.NewUserConfigService(container.UserConfigRepository)
	container.SearchService = application.NewSearchService(container.SearchRepository)

	// Initialize HTTP handlers
	container.ProjectHandler = http.NewProjectHandler(container.ProjectService)
	container.TestSuiteHandler = http.NewTestSuiteHandler(container.TestSuiteService)
	container.TestCaseHandler = http.NewTestCaseHandler(container.TestCaseService)
	container.BuildHandler = http.NewBuildHandler(container.BuildService)
	container.BuildTestCaseExecutionHandler = http.NewBuildTestCaseExecutionHandler(container.BuildTestCaseExecutionService)
	container.FailureHandler = http.NewFailureHandler(container.FailureService)
	container.UserHandler = http.NewUserHandler(container.UserService)
	container.UserConfigHandler = http.NewUserConfigHandler(container.UserConfigService)
	container.SearchHandler = http.NewSearchHandler(container.SearchService)

	container.JUnitImportService = application.NewJUnitImportService(
		container.BuildService,
		container.TestSuiteService,
		container.TestCaseService,
		application.NewBuildExecutionServiceAdapter(container.BuildTestCaseExecutionService),
	)
	container.JUnitImportHandler = http.NewJUnitImportHandler(container.JUnitImportService)

	return container
}

// GetProjectHandler returns the project handler
func (c *Container) GetProjectHandler() *http.ProjectHandler {
	return c.ProjectHandler
}

// GetProjectService returns the project service
func (c *Container) GetProjectService() domain.ProjectService {
	return c.ProjectService
}

// GetProjectRepository returns the project repository
func (c *Container) GetProjectRepository() domain.ProjectRepository {
	return c.ProjectRepository
}

// GetTestSuiteHandler returns the test suite handler
func (c *Container) GetTestSuiteHandler() *http.TestSuiteHandler {
	return c.TestSuiteHandler
}

// GetTestSuiteService returns the test suite service
func (c *Container) GetTestSuiteService() domain.TestSuiteService {
	return c.TestSuiteService
}

// GetTestSuiteRepository returns the test suite repository
func (c *Container) GetTestSuiteRepository() domain.TestSuiteRepository {
	return c.TestSuiteRepository
}

// GetTestCaseHandler returns the test case handler
func (c *Container) GetTestCaseHandler() *http.TestCaseHandler {
	return c.TestCaseHandler
}

// GetTestCaseService returns the test case service
func (c *Container) GetTestCaseService() domain.TestCaseService {
	return c.TestCaseService
}

// GetTestCaseRepository returns the test case repository
func (c *Container) GetTestCaseRepository() domain.TestCaseRepository {
	return c.TestCaseRepository
}

// GetBuildHandler returns the build handler
func (c *Container) GetBuildHandler() *http.BuildHandler {
	return c.BuildHandler
}

// GetBuildService returns the build service
func (c *Container) GetBuildService() domain.BuildService {
	return c.BuildService
}

// GetBuildRepository returns the build repository
func (c *Container) GetBuildRepository() domain.BuildRepository {
	return c.BuildRepository
}

// GetBuildTestCaseExecutionHandler returns the build test case execution handler
func (c *Container) GetBuildTestCaseExecutionHandler() *http.BuildTestCaseExecutionHandler {
	return c.BuildTestCaseExecutionHandler
}

// GetBuildTestCaseExecutionService returns the build test case execution service
func (c *Container) GetBuildTestCaseExecutionService() domain.BuildTestCaseExecutionService {
	return c.BuildTestCaseExecutionService
}

// GetBuildTestCaseExecutionRepository returns the build test case execution repository
func (c *Container) GetBuildTestCaseExecutionRepository() domain.BuildTestCaseExecutionRepository {
	return c.BuildTestCaseExecutionRepository
}

// GetFailureHandler returns the failure handler
func (c *Container) GetFailureHandler() *http.FailureHandler {
	return c.FailureHandler
}

// GetFailureService returns the failure service
func (c *Container) GetFailureService() domain.FailureService {
	return c.FailureService
}

// GetFailureRepository returns the failure repository
func (c *Container) GetFailureRepository() domain.FailureRepository {
	return c.FailureRepository
}

// GetUserHandler returns the user handler
func (c *Container) GetUserHandler() *http.UserHandler {
	return c.UserHandler
}

// GetUserService returns the user service
func (c *Container) GetUserService() domain.UserService {
	return c.UserService
}

// GetUserRepository returns the user repository
func (c *Container) GetUserRepository() domain.UserRepository {
	return c.UserRepository
}

// GetUserConfigHandler returns the user config handler
func (c *Container) GetUserConfigHandler() *http.UserConfigHandler {
	return c.UserConfigHandler
}

// GetUserConfigService returns the user config service
func (c *Container) GetUserConfigService() domain.UserConfigService {
	return c.UserConfigService
}

// GetUserConfigRepository returns the user config repository
func (c *Container) GetUserConfigRepository() domain.UserConfigRepository {
	return c.UserConfigRepository
}

// GetSearchHandler returns the search handler
func (c *Container) GetSearchHandler() *http.SearchHandler {
	return c.SearchHandler
}

// GetSearchService returns the search service
func (c *Container) GetSearchService() domain.SearchService {
	return c.SearchService
}

// GetSearchRepository returns the search repository
func (c *Container) GetSearchRepository() domain.SearchRepository {
	return c.SearchRepository
}

// GetJUnitImportHandler returns the JUnit import handler
func (c *Container) GetJUnitImportHandler() *http.JUnitImportHandler {
	return c.JUnitImportHandler
}
