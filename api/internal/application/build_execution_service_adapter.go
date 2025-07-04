package application

import (
	"context"

	"github.com/BennyEisner/test-results/internal/domain"
)

type BuildExecutionServiceAdapter struct {
	inner domain.BuildTestCaseExecutionService
}

func NewBuildExecutionServiceAdapter(inner domain.BuildTestCaseExecutionService) domain.BuildExecutionService {
	return &BuildExecutionServiceAdapter{inner: inner}
}

func (a *BuildExecutionServiceAdapter) CreateBuildExecutions(ctx context.Context, buildID int64, executions []*domain.BuildExecution) error {
	for _, exec := range executions {
		if exec == nil {
			continue
		}
		input := &domain.BuildExecutionInput{
			TestCaseID:    exec.TestCaseID,
			Status:        exec.Status,
			ExecutionTime: exec.ExecutionTime,
		}
		_, err := a.inner.CreateExecution(ctx, buildID, input)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *BuildExecutionServiceAdapter) GetBuildExecutions(ctx context.Context, buildID int64) ([]*domain.BuildExecution, error) {
	return []*domain.BuildExecution{}, nil // Not implemented
}
