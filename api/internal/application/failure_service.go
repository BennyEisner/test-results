package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// FailureService implements the FailureService interface
type FailureService struct {
	repo ports.FailureRepository
}

func NewFailureService(repo ports.FailureRepository) ports.FailureService {
	return &FailureService{repo: repo}
}

func (s *FailureService) GetFailureByID(ctx context.Context, id int64) (*models.Failure, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid failure ID")
	}

	failure, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get failure by ID %d: %w", id, err)
	}
	if failure == nil {
		return nil, domain.ErrFailureNotFound
	}
	return failure, nil
}

func (s *FailureService) GetFailureByExecutionID(ctx context.Context, executionID int64) (*models.Failure, error) {
	if executionID <= 0 {
		return nil, domain.ErrInvalidExecutionData
	}

	failure, err := s.repo.GetByExecutionID(ctx, executionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get failure by execution ID %d: %w", executionID, err)
	}
	if failure == nil {
		return nil, domain.ErrFailureNotFound
	}
	return failure, nil
}

func (s *FailureService) CreateFailure(ctx context.Context, failure *models.Failure) (*models.Failure, error) {
	if failure == nil {
		return nil, fmt.Errorf("failure cannot be nil")
	}
	if failure.BuildTestCaseExecutionID <= 0 {
		return nil, fmt.Errorf("execution ID is required")
	}

	if err := s.repo.Create(ctx, failure); err != nil {
		return nil, fmt.Errorf("failed to create failure: %w", err)
	}
	return failure, nil
}

func (s *FailureService) UpdateFailure(ctx context.Context, id int64, failure *models.Failure) (*models.Failure, error) {
	if id <= 0 || failure == nil {
		return nil, fmt.Errorf("invalid input")
	}

	updatedFailure, err := s.repo.Update(ctx, id, failure)
	if err != nil {
		return nil, fmt.Errorf("failed to update failure: %w", err)
	}
	return updatedFailure, nil
}

func (s *FailureService) DeleteFailure(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid failure ID")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete failure: %w", err)
	}
	return nil
}
