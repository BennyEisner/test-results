package application

import (
	"context"

	"github.com/BennyEisner/test-results/internal/build/domain/models"
	"github.com/BennyEisner/test-results/internal/build/domain/ports"
)

type BuildService interface {
	GetBuilds(ctx context.Context, projectID int64, suiteID *int64) ([]*models.Build, error)
	GetBuildByID(ctx context.Context, id int64) (*models.Build, error)
	CreateBuild(ctx context.Context, build *models.Build) (int64, error)
	UpdateBuild(ctx context.Context, build *models.Build) error
	DeleteBuild(ctx context.Context, id int64) error
	GetBuildDurationTrends(ctx context.Context, projectID int64, suiteID int64) ([]*models.BuildDurationTrend, error)
}

type BuildServiceImpl struct {
	repo ports.BuildRepository
}

func NewBuildService(repo ports.BuildRepository) BuildService {
	return &BuildServiceImpl{repo: repo}
}

func (s *BuildServiceImpl) GetBuilds(ctx context.Context, projectID int64, suiteID *int64) ([]*models.Build, error) {
	return s.repo.GetBuilds(ctx, projectID, suiteID)
}

func (s *BuildServiceImpl) GetBuildByID(ctx context.Context, id int64) (*models.Build, error) {
	return s.repo.GetBuildByID(ctx, id)
}

func (s *BuildServiceImpl) GetBuildDurationTrends(ctx context.Context, projectID int64, suiteID int64) ([]*models.BuildDurationTrend, error) {
	return s.repo.GetBuildDurationTrends(ctx, projectID, suiteID)
}

func (s *BuildServiceImpl) CreateBuild(ctx context.Context, build *models.Build) (int64, error) {
	return s.repo.CreateBuild(ctx, build)
}

func (s *BuildServiceImpl) UpdateBuild(ctx context.Context, build *models.Build) error {
	return s.repo.UpdateBuild(ctx, build)
}

func (s *BuildServiceImpl) DeleteBuild(ctx context.Context, id int64) error {
	return s.repo.DeleteBuild(ctx, id)
}
