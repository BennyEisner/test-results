package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/domain/models"
)

// ProjectRepository defines the interface for project data access
type ProjectRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Project, error)
	GetAll(ctx context.Context) ([]*models.Project, error)
	GetByName(ctx context.Context, name string) (*models.Project, error)
	Create(ctx context.Context, p *models.Project) error
	Update(ctx context.Context, id int64, name string) (*models.Project, error)
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int, error)
}

// BuildRepository defines the interface for build data access
type BuildRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Build, error)
	GetAllByProjectID(ctx context.Context, projectID int64) ([]*models.Build, error)
	GetAllByTestSuiteID(ctx context.Context, suiteID int64) ([]*models.Build, error)
	Create(ctx context.Context, build *models.Build) error
	Update(ctx context.Context, id int64, build *models.Build) (*models.Build, error)
	Delete(ctx context.Context, id int64) error
}

// TestSuiteRepository defines the interface for test suite data access
type TestSuiteRepository interface {
	GetByID(ctx context.Context, id int64) (*models.TestSuite, error)
	GetAllByProjectID(ctx context.Context, projectID int64) ([]*models.TestSuite, error)
	GetByName(ctx context.Context, projectID int64, name string) (*models.TestSuite, error)
	Create(ctx context.Context, suite *models.TestSuite) error
	Update(ctx context.Context, id int64, name string) (*models.TestSuite, error)
	Delete(ctx context.Context, id int64) error
}

// TestCaseRepository defines the interface for test case data access
type TestCaseRepository interface {
	GetByID(ctx context.Context, id int64) (*models.TestCase, error)
	GetAllBySuiteID(ctx context.Context, suiteID int64) ([]*models.TestCase, error)
	GetByName(ctx context.Context, suiteID int64, name string) (*models.TestCase, error)
	Create(ctx context.Context, tc *models.TestCase) error
	Update(ctx context.Context, id int64, name, classname string) (*models.TestCase, error)
	Delete(ctx context.Context, id int64) error
}

// BuildExecutionRepository defines the interface for build execution data access
type BuildExecutionRepository interface {
	GetByBuildID(ctx context.Context, buildID int64) ([]*models.BuildExecution, error)
	Create(ctx context.Context, execution *models.BuildExecution) error
	CreateBatch(ctx context.Context, executions []*models.BuildExecution) error
}

// BuildTestCaseExecutionRepository defines the interface for build test case execution data access
type BuildTestCaseExecutionRepository interface {
	GetByID(ctx context.Context, id int64) (*models.BuildTestCaseExecution, error)
	GetAllByBuildID(ctx context.Context, buildID int64) ([]*models.BuildExecutionDetail, error)
	Create(ctx context.Context, execution *models.BuildTestCaseExecution) error
	Update(ctx context.Context, id int64, execution *models.BuildTestCaseExecution) (*models.BuildTestCaseExecution, error)
	Delete(ctx context.Context, id int64) error
}

// FailureRepository defines the interface for failure data access
type FailureRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Failure, error)
	GetByExecutionID(ctx context.Context, executionID int64) (*models.Failure, error)
	Create(ctx context.Context, failure *models.Failure) error
	Update(ctx context.Context, id int64, failure *models.Failure) (*models.Failure, error)
	Delete(ctx context.Context, id int64) error
}

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(ctx context.Context, id int) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, id int, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id int) error
}

// UserConfigRepository defines the interface for user config data access
type UserConfigRepository interface {
	GetByUserID(ctx context.Context, userID int) (*models.UserConfig, error)
	Create(ctx context.Context, config *models.UserConfig) error
	Update(ctx context.Context, userID int, config *models.UserConfig) (*models.UserConfig, error)
	Delete(ctx context.Context, userID int) error
}

// SearchRepository defines the interface for search data access
type SearchRepository interface {
	Search(ctx context.Context, query string) ([]*models.SearchResult, error)
}
