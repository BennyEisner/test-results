package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/build/domain/models"
)

// BuildRepository defines the interface for build data access
type BuildRepository interface {
	GetByID(ctx context.Context, id int64) (*models.Build, error)
	GetAllByProjectID(ctx context.Context, projectID int64) ([]*models.Build, error)
	GetAllByTestSuiteID(ctx context.Context, suiteID int64) ([]*models.Build, error)
	Create(ctx context.Context, build *models.Build) error
	Update(ctx context.Context, id int64, build *models.Build) (*models.Build, error)
	Delete(ctx context.Context, id int64) error
}

// BuildService defines the interface for build business logic
type BuildService interface {
	GetBuild(ctx context.Context, id int64) (*models.Build, error)
	GetBuildsByProject(ctx context.Context, projectID int64) ([]*models.Build, error)
	GetBuildsByTestSuite(ctx context.Context, suiteID int64) ([]*models.Build, error)
	CreateBuild(ctx context.Context, projectID int64, suiteID int64, name string) (*models.Build, error)
	UpdateBuild(ctx context.Context, id int64, name string) (*models.Build, error)
	DeleteBuild(ctx context.Context, id int64) error
}
