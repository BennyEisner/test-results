package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
)

type BuildService struct {
	repo domain.BuildRepository
}

func NewBuildService(repo domain.BuildRepository) domain.BuildService {
	return &BuildService{repo: repo}
}

func (s *BuildService) GetBuildByID(ctx context.Context, id int64) (*domain.Build, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
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

func (s *BuildService) GetBuildsByProjectID(ctx context.Context, projectID int64) ([]*domain.Build, error) {
	if projectID <= 0 {
		return nil, domain.ErrInvalidInput
	}
	builds, err := s.repo.GetAllByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds by project ID %d: %w", projectID, err)
	}
	return builds, nil
}

func (s *BuildService) GetBuildsByTestSuiteID(ctx context.Context, suiteID int64) ([]*domain.Build, error) {
	if suiteID <= 0 {
		return nil, domain.ErrInvalidInput
	}
	builds, err := s.repo.GetAllByTestSuiteID(ctx, suiteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds by test suite ID %d: %w", suiteID, err)
	}
	return builds, nil
}

func (s *BuildService) CreateBuild(ctx context.Context, build *domain.Build) (*domain.Build, error) {
	if build == nil {
		return nil, domain.ErrInvalidInput
	}
	if build.TestSuiteID <= 0 || build.ProjectID <= 0 || build.BuildNumber == "" {
		return nil, domain.ErrInvalidInput
	}
	if err := s.repo.Create(ctx, build); err != nil {
		return nil, fmt.Errorf("failed to create build: %w", err)
	}
	return build, nil
}

func (s *BuildService) UpdateBuild(ctx context.Context, id int64, build *domain.Build) (*domain.Build, error) {
	if id <= 0 || build == nil {
		return nil, domain.ErrInvalidInput
	}
	updatedBuild, err := s.repo.Update(ctx, id, build)
	if err != nil {
		return nil, fmt.Errorf("failed to update build: %w", err)
	}
	return updatedBuild, nil
}

func (s *BuildService) DeleteBuild(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete build: %w", err)
	}
	return nil
}
