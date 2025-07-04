package application

import (
	"context"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// BuildExecutionServiceAdapter adapts BuildTestCaseExecutionService to BuildExecutionService
type BuildExecutionServiceAdapter struct {
	buildTestCaseExecutionService ports.BuildTestCaseExecutionService
}

// NewBuildExecutionServiceAdapter creates a new adapter
func NewBuildExecutionServiceAdapter(buildTestCaseExecutionService ports.BuildTestCaseExecutionService) ports.BuildExecutionService {
	return &BuildExecutionServiceAdapter{
		buildTestCaseExecutionService: buildTestCaseExecutionService,
	}
}

// GetBuildExecutions implements BuildExecutionService
func (a *BuildExecutionServiceAdapter) GetBuildExecutions(ctx context.Context, buildID int64) ([]*models.BuildExecution, error) {
	// This is a stub implementation - the adapter doesn't need to implement this for JUnitImportService
	return nil, nil
}

// CreateBuildExecutions implements BuildExecutionService
func (a *BuildExecutionServiceAdapter) CreateBuildExecutions(ctx context.Context, buildID int64, executions []*models.BuildExecution) error {
	for _, execution := range executions {
		if execution == nil {
			continue
		}

		input := &models.BuildExecutionInput{
			TestCaseID:    execution.TestCaseID,
			Status:        execution.Status,
			ExecutionTime: execution.ExecutionTime,
		}

		_, err := a.buildTestCaseExecutionService.CreateExecution(ctx, buildID, input)
		if err != nil {
			return err
		}
	}
	return nil
}
