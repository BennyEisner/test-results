package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/test_case/domain/models"
)

// TestCaseRepository defines the interface for test case data access
type TestCaseRepository interface {
	GetByID(ctx context.Context, id int64) (*models.TestCase, error)
	GetAllBySuiteID(ctx context.Context, suiteID int64) ([]*models.TestCase, error)
	GetByName(ctx context.Context, suiteID int64, name string) (*models.TestCase, error)
	Create(ctx context.Context, tc *models.TestCase) error
	Update(ctx context.Context, id int64, name, classname string) (*models.TestCase, error)
	Delete(ctx context.Context, id int64) error
}

// TestCaseService defines the interface for test case business logic
type TestCaseService interface {
	GetTestCase(ctx context.Context, id int64) (*models.TestCase, error)
	GetTestCasesBySuite(ctx context.Context, suiteID int64) ([]*models.TestCase, error)
	CreateTestCase(ctx context.Context, suiteID int64, name, classname string) (*models.TestCase, error)
	UpdateTestCase(ctx context.Context, id int64, name, classname string) (*models.TestCase, error)
	DeleteTestCase(ctx context.Context, id int64) error
}
