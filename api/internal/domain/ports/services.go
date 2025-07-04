package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/domain/models"
)

// ProjectService defines the business logic for project operations
type ProjectService interface {
	GetProjectByID(ctx context.Context, id int64) (*models.Project, error)
	GetAllProjects(ctx context.Context) ([]*models.Project, error)
	CreateProject(ctx context.Context, name string) (*models.Project, error)
	UpdateProject(ctx context.Context, id int64, name string) (*models.Project, error)
	DeleteProject(ctx context.Context, id int64) error
	GetProjectByName(ctx context.Context, name string) (*models.Project, error)
}

// BuildService defines the business logic for build operations
type BuildService interface {
	GetBuildByID(ctx context.Context, id int64) (*models.Build, error)
	GetBuildsByProjectID(ctx context.Context, projectID int64) ([]*models.Build, error)
	GetBuildsByTestSuiteID(ctx context.Context, suiteID int64) ([]*models.Build, error)
	CreateBuild(ctx context.Context, build *models.Build) (*models.Build, error)
	UpdateBuild(ctx context.Context, id int64, build *models.Build) (*models.Build, error)
	DeleteBuild(ctx context.Context, id int64) error
}

// TestSuiteService defines the business logic for test suite operations
type TestSuiteService interface {
	GetTestSuiteByID(ctx context.Context, id int64) (*models.TestSuite, error)
	GetTestSuitesByProjectID(ctx context.Context, projectID int64) ([]*models.TestSuite, error)
	GetTestSuiteByName(ctx context.Context, projectID int64, name string) (*models.TestSuite, error)
	CreateTestSuite(ctx context.Context, projectID int64, name string, parentID *int64) (*models.TestSuite, error)
	UpdateTestSuite(ctx context.Context, id int64, name string) (*models.TestSuite, error)
	DeleteTestSuite(ctx context.Context, id int64) error
}

// TestCaseService defines the business logic for test case operations
type TestCaseService interface {
	GetTestCaseByID(ctx context.Context, id int64) (*models.TestCase, error)
	GetTestCasesBySuiteID(ctx context.Context, suiteID int64) ([]*models.TestCase, error)
	GetTestCaseByName(ctx context.Context, suiteID int64, name string) (*models.TestCase, error)
	CreateTestCase(ctx context.Context, suiteID int64, name, classname string) (*models.TestCase, error)
	UpdateTestCase(ctx context.Context, id int64, name, classname string) (*models.TestCase, error)
	DeleteTestCase(ctx context.Context, id int64) error
}

// BuildExecutionService defines the business logic for build execution operations
type BuildExecutionService interface {
	GetBuildExecutions(ctx context.Context, buildID int64) ([]*models.BuildExecution, error)
	CreateBuildExecutions(ctx context.Context, buildID int64, executions []*models.BuildExecution) error
}

// BuildTestCaseExecutionService defines the business logic for build test case execution operations
type BuildTestCaseExecutionService interface {
	GetExecutionByID(ctx context.Context, id int64) (*models.BuildTestCaseExecution, error)
	GetExecutionsByBuildID(ctx context.Context, buildID int64) ([]*models.BuildExecutionDetail, error)
	CreateExecution(ctx context.Context, buildID int64, input *models.BuildExecutionInput) (*models.BuildTestCaseExecution, error)
	UpdateExecution(ctx context.Context, id int64, execution *models.BuildTestCaseExecution) (*models.BuildTestCaseExecution, error)
	DeleteExecution(ctx context.Context, id int64) error
}

// FailureService defines the business logic for failure operations
type FailureService interface {
	GetFailureByID(ctx context.Context, id int64) (*models.Failure, error)
	GetFailureByExecutionID(ctx context.Context, executionID int64) (*models.Failure, error)
	CreateFailure(ctx context.Context, failure *models.Failure) (*models.Failure, error)
	UpdateFailure(ctx context.Context, id int64, failure *models.Failure) (*models.Failure, error)
	DeleteFailure(ctx context.Context, id int64) error
}

// UserService defines the business logic for user operations
type UserService interface {
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	CreateUser(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, id int, username string) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
}

// UserConfigService defines the business logic for user config operations
type UserConfigService interface {
	GetUserConfig(ctx context.Context, userID int) (*models.UserConfig, error)
	CreateUserConfig(ctx context.Context, userID int, layouts, activeLayoutID string) (*models.UserConfig, error)
	UpdateUserConfig(ctx context.Context, userID int, layouts, activeLayoutID string) (*models.UserConfig, error)
	DeleteUserConfig(ctx context.Context, userID int) error
}

// JUnitImportService defines the business logic for JUnit import operations
type JUnitImportService interface {
	ProcessJUnitData(ctx context.Context, projectID int64, suiteID int64, junitData *models.JUnitTestSuites) (*models.Build, error)
}

// SearchService defines the business logic for search operations
type SearchService interface {
	Search(ctx context.Context, query string) ([]*models.SearchResult, error)
}
