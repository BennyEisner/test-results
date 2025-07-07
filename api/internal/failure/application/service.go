package application

import (
	"context"
	"fmt"
	"time"

	"github.com/BennyEisner/test-results/internal/failure/domain/errors"
	"github.com/BennyEisner/test-results/internal/failure/domain/models"
	"github.com/BennyEisner/test-results/internal/failure/domain/ports"
)

// FailureService implements the FailureService interface
type FailureService struct {
	repo ports.FailureRepository
}

func NewFailureService(repo ports.FailureRepository) ports.FailureService {
	return &FailureService{repo: repo}
}

func (s *FailureService) GetFailure(ctx context.Context, id int64) (*models.Failure, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid failure ID")
	}

	failure, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get failure by ID %d: %w", id, err)
	}
	if failure == nil {
		return nil, errors.ErrFailureNotFound
	}
	return failure, nil
}

func (s *FailureService) GetFailureByExecution(ctx context.Context, executionID int64) (*models.Failure, error) {
	if executionID <= 0 {
		return nil, fmt.Errorf("invalid execution ID")
	}

	failure, err := s.repo.GetByExecutionID(ctx, executionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get failure by execution ID %d: %w", executionID, err)
	}
	if failure == nil {
		return nil, errors.ErrFailureNotFound
	}
	return failure, nil
}

func (s *FailureService) CreateFailure(ctx context.Context, executionID int64, message, failureType, details string) (*models.Failure, error) {
	if executionID <= 0 {
		return nil, fmt.Errorf("execution ID is required")
	}
	if message == "" {
		return nil, fmt.Errorf("failure message is required")
	}

	failure := &models.Failure{
		ExecutionID: executionID,
		Message:     message,
		Type:        failureType,
		Details:     details,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, failure); err != nil {
		return nil, fmt.Errorf("failed to create failure: %w", err)
	}
	return failure, nil
}

func (s *FailureService) UpdateFailure(ctx context.Context, id int64, message, failureType, details string) (*models.Failure, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid failure ID")
	}

	failure := &models.Failure{
		Message: message,
		Type:    failureType,
		Details: details,
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
