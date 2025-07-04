package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// BuildTestCaseExecutionService implements the BuildTestCaseExecutionService interface
type BuildTestCaseExecutionService struct {
	repo ports.BuildTestCaseExecutionRepository
}

func NewBuildTestCaseExecutionService(repo ports.BuildTestCaseExecutionRepository) ports.BuildTestCaseExecutionService {
	return &BuildTestCaseExecutionService{repo: repo}
}

func (s *BuildTestCaseExecutionService) GetExecutionByID(ctx context.Context, id int64) (*models.BuildTestCaseExecution, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidExecutionData
	}

	execution, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution by ID %d: %w", id, err)
	}
	if execution == nil {
		return nil, domain.ErrBuildExecutionNotFound
	}
	return execution, nil
}

func (s *BuildTestCaseExecutionService) GetExecutionsByBuildID(ctx context.Context, buildID int64) ([]*models.BuildExecutionDetail, error) {
	if buildID <= 0 {
		return nil, domain.ErrInvalidBuildData
	}
	executions, err := s.repo.GetAllByBuildID(ctx, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get executions by build ID %d: %w", buildID, err)
	}
	return executions, nil
}

func (s *BuildTestCaseExecutionService) CreateExecution(ctx context.Context, buildID int64, input *models.BuildExecutionInput) (*models.BuildTestCaseExecution, error) {
	if buildID <= 0 || input == nil {
		return nil, domain.ErrInvalidExecutionData
	}
	if input.TestCaseID <= 0 || input.Status == "" {
		return nil, domain.ErrInvalidExecutionData
	}

	execution := &models.BuildTestCaseExecution{
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

func (s *BuildTestCaseExecutionService) UpdateExecution(ctx context.Context, id int64, execution *models.BuildTestCaseExecution) (*models.BuildTestCaseExecution, error) {
	if id <= 0 || execution == nil {
		return nil, domain.ErrInvalidExecutionData
	}

	updatedExecution, err := s.repo.Update(ctx, id, execution)
	if err != nil {
		return nil, fmt.Errorf("failed to update execution: %w", err)
	}
	return updatedExecution, nil
}

func (s *BuildTestCaseExecutionService) DeleteExecution(ctx context.Context, id int64) error {
	if id <= 0 {
		return domain.ErrInvalidExecutionData
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete execution: %w", err)
	}
	return nil
}

// Ensure BuildTestCaseExecutionService implements BuildExecutionService
var _ ports.BuildExecutionService = (*BuildTestCaseExecutionService)(nil)

func (s *BuildTestCaseExecutionService) CreateBuildExecutions(ctx context.Context, buildID int64, executions []*models.BuildExecution) error {
	for _, exec := range executions {
		if exec == nil {
			continue
		}
		input := &models.BuildExecutionInput{
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

func (s *BuildTestCaseExecutionService) GetBuildExecutions(ctx context.Context, buildID int64) ([]*models.BuildExecution, error) {
	return nil, nil // Not used in JUnitImportService, implement as stub or map from GetExecutionsByBuildID if needed
}
