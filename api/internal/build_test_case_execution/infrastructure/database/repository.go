package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/build_test_case_execution/domain/models"
	"github.com/BennyEisner/test-results/internal/build_test_case_execution/domain/ports"
)

// SQLBuildTestCaseExecutionRepository implements BuildTestCaseExecutionRepository
type SQLBuildTestCaseExecutionRepository struct {
	db *sql.DB
}

// NewSQLBuildTestCaseExecutionRepository creates a new SQL repository
func NewSQLBuildTestCaseExecutionRepository(db *sql.DB) ports.BuildTestCaseExecutionRepository {
	return &SQLBuildTestCaseExecutionRepository{db: db}
}

// GetByID retrieves a build test case execution by ID
func (r *SQLBuildTestCaseExecutionRepository) GetByID(ctx context.Context, id int64) (*models.BuildTestCaseExecution, error) {
	query := `SELECT id, build_id, test_case_id, status, execution_time, created_at 
			  FROM build_test_case_executions WHERE id = $1`

	var execution models.BuildTestCaseExecution
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&execution.ID,
		&execution.BuildID,
		&execution.TestCaseID,
		&execution.Status,
		&execution.ExecutionTime,
		&execution.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get execution by ID: %w", err)
	}

	return &execution, nil
}

// GetAllByBuildID retrieves all build test case executions for a build
func (r *SQLBuildTestCaseExecutionRepository) GetAllByBuildID(ctx context.Context, buildID int64) ([]*models.BuildExecutionDetail, error) {
	query := `SELECT e.id, e.build_id, e.test_case_id, tc.name, tc.classname, 
			  e.status, e.execution_time, e.created_at
			  FROM build_test_case_executions e
			  JOIN test_cases tc ON e.test_case_id = tc.id
			  WHERE e.build_id = $1`

	rows, err := r.db.QueryContext(ctx, query, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get executions by build ID: %w", err)
	}
	defer rows.Close()

	var executions []*models.BuildExecutionDetail
	for rows.Next() {
		var execution models.BuildExecutionDetail
		err := rows.Scan(
			&execution.ExecutionID,
			&execution.BuildID,
			&execution.TestCaseID,
			&execution.TestCaseName,
			&execution.ClassName,
			&execution.Status,
			&execution.ExecutionTime,
			&execution.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan execution: %w", err)
		}
		executions = append(executions, &execution)
	}

	return executions, nil
}

// Create creates a new build test case execution
func (r *SQLBuildTestCaseExecutionRepository) Create(ctx context.Context, execution *models.BuildTestCaseExecution) error {
	query := `INSERT INTO build_test_case_executions (build_id, test_case_id, status, execution_time)
			  VALUES ($1, $2, $3, $4) RETURNING id, created_at`

	err := r.db.QueryRowContext(ctx, query,
		execution.BuildID,
		execution.TestCaseID,
		execution.Status,
		execution.ExecutionTime,
	).Scan(&execution.ID, &execution.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create execution: %w", err)
	}

	return nil
}

// Update updates an existing build test case execution
func (r *SQLBuildTestCaseExecutionRepository) Update(ctx context.Context, id int64, execution *models.BuildTestCaseExecution) (*models.BuildTestCaseExecution, error) {
	query := `UPDATE build_test_case_executions 
			  SET build_id = $1, test_case_id = $2, status = $3, execution_time = $4
			  WHERE id = $5 RETURNING id, build_id, test_case_id, status, execution_time, created_at`

	var updatedExecution models.BuildTestCaseExecution
	err := r.db.QueryRowContext(ctx, query,
		execution.BuildID,
		execution.TestCaseID,
		execution.Status,
		execution.ExecutionTime,
		id,
	).Scan(
		&updatedExecution.ID,
		&updatedExecution.BuildID,
		&updatedExecution.TestCaseID,
		&updatedExecution.Status,
		&updatedExecution.ExecutionTime,
		&updatedExecution.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update execution: %w", err)
	}

	return &updatedExecution, nil
}

// Delete deletes a build test case execution
func (r *SQLBuildTestCaseExecutionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM build_test_case_executions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete execution: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("execution not found")
	}

	return nil
}
