package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/build_test_case_execution/domain/models"
)

// BuildTestCaseExecutionRepository defines the interface for build test case execution data access
type BuildTestCaseExecutionRepository interface {
	GetByID(ctx context.Context, id int64) (*models.BuildTestCaseExecution, error)
	GetAllByBuildID(ctx context.Context, buildID int64) ([]*models.BuildExecutionDetail, error)
	Create(ctx context.Context, execution *models.BuildTestCaseExecution) error
	Update(ctx context.Context, id int64, execution *models.BuildTestCaseExecution) (*models.BuildTestCaseExecution, error)
	Delete(ctx context.Context, id int64) error
}

// BuildTestCaseExecutionService defines the interface for build test case execution business logic
type BuildTestCaseExecutionService interface {
	GetExecutionByID(ctx context.Context, id int64) (*models.BuildTestCaseExecution, error)
	GetExecutionsByBuildID(ctx context.Context, buildID int64) ([]*models.BuildExecutionDetail, error)
	CreateExecution(ctx context.Context, buildID int64, input *models.BuildExecutionInput) (*models.BuildTestCaseExecution, error)
	UpdateExecution(ctx context.Context, id int64, execution *models.BuildTestCaseExecution) (*models.BuildTestCaseExecution, error)
	DeleteExecution(ctx context.Context, id int64) error
}

// BuildExecutionService defines the interface for build execution business logic (adapter pattern)
type BuildExecutionService interface {
	CreateBuildExecutions(ctx context.Context, buildID int64, executions []*models.BuildExecution) error
	GetBuildExecutions(ctx context.Context, buildID int64) ([]*models.BuildExecution, error)
}
