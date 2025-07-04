package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
)

type BuildTestCaseExecutionService struct {
	repo domain.BuildTestCaseExecutionRepository
}

func NewBuildTestCaseExecutionService(repo domain.BuildTestCaseExecutionRepository) domain.BuildTestCaseExecutionService {
	return &BuildTestCaseExecutionService{repo: repo}
}

func (s *BuildTestCaseExecutionService) GetExecutionByID(ctx context.Context, id int64) (*domain.BuildTestCaseExecution, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
	}
	execution, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution by ID %d: %w", id, err)
	}
	if execution == nil {
		return nil, domain.ErrExecutionNotFound
	}
	return execution, nil
}

func (s *BuildTestCaseExecutionService) GetExecutionsByBuildID(ctx context.Context, buildID int64) ([]*domain.BuildExecutionDetail, error) {
	if buildID <= 0 {
		return nil, domain.ErrInvalidInput
	}
	executions, err := s.repo.GetAllByBuildID(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get executions by build ID %d: %w", buildID, err)
	}
	return executions, nil
}

func (s *BuildTestCaseExecutionService) CreateExecution(ctx context.Context, buildID int64, input *domain.BuildExecutionInput) (*domain.BuildTestCaseExecution, error) {
	if buildID <= 0 || input == nil {
		return nil, domain.ErrInvalidInput
	}
	if input.TestCaseID <= 0 || input.Status == "" {
		return nil, domain.ErrInvalidInput
	}
	execution := &domain.BuildTestCaseExecution{
		BuildID:       buildID,
		TestCaseID:    input.TestCaseID,
		Status:        input.Status,
		ExecutionTime: input.ExecutionTime,
	}
	if err := s.repo.Create(ctx, execution); err != nil {
		return nil, fmt.Errorf("failed to create execution: %w", err)
	}
	return execution, nil
}

func (s *BuildTestCaseExecutionService) UpdateExecution(ctx context.Context, id int64, execution *domain.BuildTestCaseExecution) (*domain.BuildTestCaseExecution, error) {
	if id <= 0 || execution == nil {
		return nil, domain.ErrInvalidInput
	}
	updatedExecution, err := s.repo.Update(ctx, id, execution)
	if err != nil {
		return nil, fmt.Errorf("failed to update execution: %w", err)
	}
	return updatedExecution, nil
}

func (s *BuildTestCaseExecutionService) DeleteExecution(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete execution: %w", err)
	}
	return nil
}

// Ensure BuildTestCaseExecutionService implements BuildExecutionService
var _ domain.BuildExecutionService = (*BuildTestCaseExecutionService)(nil)

func (s *BuildTestCaseExecutionService) CreateBuildExecutions(ctx context.Context, buildID int64, executions []*domain.BuildExecution) error {
	for _, exec := range executions {
		if exec == nil {
			continue
		}
		input := &domain.BuildExecutionInput{
			TestCaseID:    exec.TestCaseID,
			Status:        exec.Status,
			ExecutionTime: exec.ExecutionTime,
		}
		_, err := s.CreateExecution(ctx, buildID, input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *BuildTestCaseExecutionService) GetBuildExecutions(ctx context.Context, buildID int64) ([]*domain.BuildExecution, error) {
	return nil, nil // Not used in JUnitImportService, implement as stub or map from GetExecutionsByBuildID if needed
}
