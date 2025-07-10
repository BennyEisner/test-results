package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/build/domain"
	"github.com/BennyEisner/test-results/internal/build/domain/models"
	"github.com/BennyEisner/test-results/internal/build/domain/ports"
)

// BuildService implements the BuildService interface
type BuildService struct {
	repo ports.BuildRepository
}

func NewBuildService(repo ports.BuildRepository) ports.BuildService {
	return &BuildService{repo: repo}
}

func (s *BuildService) GetBuild(ctx context.Context, id int64) (*models.Build, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidBuildData
	}

	build, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get build by ID %d: %w", id, err)
	}
	if build == nil {
		return nil, domain.ErrBuildNotFound
	}
	return build, nil
}

func (s *BuildService) GetBuildsByProject(ctx context.Context, projectID int64) ([]*models.Build, error) {
	if projectID <= 0 {
		return nil, domain.ErrInvalidProjectName
	}
	builds, err := s.repo.GetAllByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds by project ID %d: %w", projectID, err)
	}
	return builds, nil
}

func (s *BuildService) GetBuildsByTestSuite(ctx context.Context, suiteID int64) ([]*models.Build, error) {
	if suiteID <= 0 {
		return nil, domain.ErrInvalidTestSuiteName
	}
	builds, err := s.repo.GetAllByTestSuiteID(ctx, suiteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds by test suite ID %d: %w", suiteID, err)
	}
	return builds, nil
}

func (s *BuildService) CreateBuild(ctx context.Context, projectID int64, suiteID int64, name string) (*models.Build, error) {
	if projectID <= 0 || suiteID <= 0 || name == "" {
		return nil, domain.ErrInvalidBuildData
	}

	build := &models.Build{
		ProjectID:   projectID,
		TestSuiteID: suiteID,
		BuildNumber: name,
	}

	if err := s.repo.Create(ctx, build); err != nil {
		return nil, fmt.Errorf("failed to create build: %w", err)
	}
	return build, nil
}

func (s *BuildService) UpdateBuild(ctx context.Context, id int64, name string) (*models.Build, error) {
	if id <= 0 || name == "" {
		return nil, domain.ErrInvalidBuildData
	}

	build := &models.Build{
		BuildNumber: name,
	}

	updatedBuild, err := s.repo.Update(ctx, id, build)
	if err != nil {
		return nil, fmt.Errorf("failed to update build: %w", err)
	}
	return updatedBuild, nil
}

func (s *BuildService) DeleteBuild(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidBuildData
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete build: %w", err)
	}
	return nil
}
