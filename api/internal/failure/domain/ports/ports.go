package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/failure/domain/models"
)

// FailureRepository defines the interface for failure data access
type FailureRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Failure, error)
	GetByExecutionID(ctx context.Context, executionID int64) (*models.Failure, error)
	Create(ctx context.Context, failure *models.Failure) error
	Update(ctx context.Context, id int64, failure *models.Failure) (*models.Failure, error)
	Delete(ctx context.Context, id int64) error
}

// FailureService defines the interface for failure business logic
type FailureService interface {
	GetFailure(ctx context.Context, id int64) (*models.Failure, error)
	GetFailureByExecution(ctx context.Context, executionID int64) (*models.Failure, error)
	CreateFailure(ctx context.Context, executionID int64, message, failureType, details string) (*models.Failure, error)
	UpdateFailure(ctx context.Context, id int64, message, failureType, details string) (*models.Failure, error)
	DeleteFailure(ctx context.Context, id int64) error
}
