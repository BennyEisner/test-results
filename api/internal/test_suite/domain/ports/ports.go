package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/test_suite/domain/models"
)

// TestSuiteRepository defines the interface for test suite data access
type TestSuiteRepository interface {
	GetByID(ctx context.Context, id int64) (*models.TestSuite, error)
	GetAllByProjectID(ctx context.Context, projectID int64) ([]*models.TestSuite, error)
	GetByName(ctx context.Context, projectID int64, name string) (*models.TestSuite, error)
	Create(ctx context.Context, suite *models.TestSuite) error
	Update(ctx context.Context, id int64, name string) (*models.TestSuite, error)
	Delete(ctx context.Context, id int64) error
}

// TestSuiteService defines the interface for test suite business logic
type TestSuiteService interface {
	GetTestSuite(ctx context.Context, id int64) (*models.TestSuite, error)
	GetTestSuitesByProject(ctx context.Context, projectID int64) ([]*models.TestSuite, error)
	CreateTestSuite(ctx context.Context, projectID int64, name string, parentID *int64, time float64) (*models.TestSuite, error)
	UpdateTestSuite(ctx context.Context, id int64, name string) (*models.TestSuite, error)
	DeleteTestSuite(ctx context.Context, id int64) error
}
