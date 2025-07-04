package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
)

type FailureService struct {
	repo domain.FailureRepository
}

func NewFailureService(repo domain.FailureRepository) domain.FailureService {
	return &FailureService{repo: repo}
}

func (s *FailureService) GetFailureByID(ctx context.Context, id int64) (*domain.Failure, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
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

func (s *FailureService) GetFailureByExecutionID(ctx context.Context, executionID int64) (*domain.Failure, error) {
	if executionID <= 0 {
		return nil, domain.ErrInvalidInput
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

func (s *FailureService) CreateFailure(ctx context.Context, failure *domain.Failure) (*domain.Failure, error) {
	if failure == nil {
		return nil, domain.ErrInvalidInput
	}
	if failure.BuildTestCaseExecutionID <= 0 {
		return nil, domain.ErrInvalidInput
	}
	if err := s.repo.Create(ctx, failure); err != nil {
		return nil, fmt.Errorf("failed to create failure: %w", err)
	}
	return failure, nil
}

func (s *FailureService) UpdateFailure(ctx context.Context, id int64, failure *domain.Failure) (*domain.Failure, error) {
	if id <= 0 || failure == nil {
		return nil, domain.ErrInvalidInput
	}
	updatedFailure, err := s.repo.Update(ctx, id, failure)
	if err != nil {
		return nil, fmt.Errorf("failed to update failure: %w", err)
	}
	return updatedFailure, nil
}

func (s *FailureService) DeleteFailure(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete failure: %w", err)
	}
	return nil
}
