package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// SQLBuildTestCaseExecutionRepository implements the BuildTestCaseExecutionRepository interface
type SQLBuildTestCaseExecutionRepository struct {
	db *sql.DB
}

// NewSQLBuildTestCaseExecutionRepository creates a new SQL build test case execution repository
func NewSQLBuildTestCaseExecutionRepository(db *sql.DB) ports.BuildTestCaseExecutionRepository {
	return &SQLBuildTestCaseExecutionRepository{db: db}
}

// GetByID retrieves a build test case execution by its ID
func (r *SQLBuildTestCaseExecutionRepository) GetByID(ctx context.Context, id int64) (*models.BuildTestCaseExecution, error) {
	query := `SELECT id, build_id, test_case_id, status, execution_time, created_at FROM build_test_case_executions WHERE id = $1`

	var execution models.BuildTestCaseExecution

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&execution.ID, &execution.BuildID, &execution.TestCaseID, &execution.Status, &execution.ExecutionTime, &execution.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get build test case execution by ID: %w", err)
	}

	return &execution, nil
}

// GetAllByBuildID retrieves all build test case executions for a build with details
func (r *SQLBuildTestCaseExecutionRepository) GetAllByBuildID(ctx context.Context, buildID int64) ([]*models.BuildExecutionDetail, error) {
	query := `
		SELECT 
			bte.id, bte.build_id, bte.test_case_id, tc.name, tc.classname,
			bte.status, bte.execution_time, bte.created_at,
			f.id, f.message, f.type, f.details
		FROM build_test_case_executions bte
		JOIN test_cases tc ON bte.test_case_id = tc.id
		LEFT JOIN failures f ON bte.id = f.execution_id
		WHERE bte.build_id = $1
		ORDER BY bte.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get build test case executions by build ID: %w", err)
	}
	defer rows.Close()

	var details []*models.BuildExecutionDetail
	for rows.Next() {
		var detail models.BuildExecutionDetail
		var failureID sql.NullInt64
		var failureMessage, failureType, failureDetails sql.NullString

		err := rows.Scan(
			&detail.ExecutionID, &detail.BuildID, &detail.TestCaseID, &detail.TestCaseName, &detail.ClassName,
			&detail.Status, &detail.ExecutionTime, &detail.CreatedAt,
			&failureID, &failureMessage, &failureType, &failureDetails,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan build execution detail: %w", err)
		}

		if failureID.Valid {
			var message, failureTypeStr, details *string
			if failureMessage.Valid {
				message = &failureMessage.String
			}
			if failureType.Valid {
				failureTypeStr = &failureType.String
			}
			if failureDetails.Valid {
				details = &failureDetails.String
			}

			detail.Failure = &models.Failure{
				ID:      failureID.Int64,
				Message: message,
				Type:    failureTypeStr,
				Details: details,
			}
		}

		details = append(details, &detail)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating build execution details: %w", err)
	}

	return details, nil
}

// Create creates a new build test case execution
func (r *SQLBuildTestCaseExecutionRepository) Create(ctx context.Context, execution *models.BuildTestCaseExecution) error {
	query := `INSERT INTO build_test_case_executions (build_id, test_case_id, status, execution_time, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		execution.BuildID, execution.TestCaseID, execution.Status, execution.ExecutionTime, execution.CreatedAt,
	).Scan(&execution.ID)
	if err != nil {
		return fmt.Errorf("failed to create build test case execution: %w", err)
	}

	return nil
}

// Update updates an existing build test case execution
func (r *SQLBuildTestCaseExecutionRepository) Update(ctx context.Context, id int64, execution *models.BuildTestCaseExecution) (*models.BuildTestCaseExecution, error) {
	query := `UPDATE build_test_case_executions SET status = $1, execution_time = $2 WHERE id = $3 RETURNING id, build_id, test_case_id, status, execution_time, created_at`

	var updatedExecution models.BuildTestCaseExecution

	err := r.db.QueryRowContext(ctx, query, execution.Status, execution.ExecutionTime, id).Scan(
		&updatedExecution.ID, &updatedExecution.BuildID, &updatedExecution.TestCaseID, &updatedExecution.Status, &updatedExecution.ExecutionTime, &updatedExecution.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update build test case execution: %w", err)
	}

	return &updatedExecution, nil
}

// Delete deletes a build test case execution by its ID
func (r *SQLBuildTestCaseExecutionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM build_test_case_executions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete build test case execution: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("build test case execution not found")
	}

	return nil
}
