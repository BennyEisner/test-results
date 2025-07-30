package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/build/domain/models"
)

type BuildRepository interface {
	GetBuilds(ctx context.Context, projectID int64, suiteID *int64) ([]*models.Build, error)
	GetBuildByID(ctx context.Context, id int64) (*models.Build, error)
	CreateBuild(ctx context.Context, build *models.Build) (int64, error)
	UpdateBuild(ctx context.Context, build *models.Build) error
	DeleteBuild(ctx context.Context, id int64) error
	GetBuildDurationTrends(ctx context.Context, projectID int64, suiteID int64) ([]*models.BuildDurationTrend, error)
	GetLatestBuildStatus(ctx context.Context, projectID int64) (string, error)
	GetLatestBuilds(ctx context.Context, projectID int64, limit int) ([]*models.Build, error)
}
