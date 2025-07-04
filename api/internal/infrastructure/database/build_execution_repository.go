package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// SQLBuildExecutionRepository implements the BuildExecutionRepository interface
type SQLBuildExecutionRepository struct {
	db *sql.DB
}

// NewSQLBuildExecutionRepository creates a new SQL build execution repository
func NewSQLBuildExecutionRepository(db *sql.DB) ports.BuildExecutionRepository {
	return &SQLBuildExecutionRepository{db: db}
}

// GetByBuildID retrieves all build executions for a build
func (r *SQLBuildExecutionRepository) GetByBuildID(ctx context.Context, buildID int64) ([]*models.BuildExecution, error) {
	query := `SELECT id, build_id, test_case_id, status, execution_time, created_at FROM build_executions WHERE build_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get build executions by build ID: %w", err)
	}
	defer rows.Close()

	var executions []*models.BuildExecution
	for rows.Next() {
		var execution models.BuildExecution

		err := rows.Scan(
			&execution.ID, &execution.BuildID, &execution.TestCaseID, &execution.Status, &execution.ExecutionTime, &execution.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan build execution: %w", err)
		}

		executions = append(executions, &execution)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating build executions: %w", err)
	}

	return executions, nil
}

// Create creates a new build execution
func (r *SQLBuildExecutionRepository) Create(ctx context.Context, execution *models.BuildExecution) error {
	query := `INSERT INTO build_executions (build_id, test_case_id, status, execution_time, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		execution.BuildID, execution.TestCaseID, execution.Status, execution.ExecutionTime, execution.CreatedAt,
	).Scan(&execution.ID)
	if err != nil {
		return fmt.Errorf("failed to create build execution: %w", err)
	}

	return nil
}

// CreateBatch creates multiple build executions in a single transaction
func (r *SQLBuildExecutionRepository) CreateBatch(ctx context.Context, executions []*models.BuildExecution) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			// Log the rollback error but don't return it since we're in a defer
			fmt.Printf("failed to rollback transaction: %v\n", err)
		}
	}()

	query := `INSERT INTO build_executions (build_id, test_case_id, status, execution_time, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	for _, execution := range executions {
		err := tx.QueryRowContext(ctx, query,
			execution.BuildID, execution.TestCaseID, execution.Status, execution.ExecutionTime, execution.CreatedAt,
		).Scan(&execution.ID)
		if err != nil {
			return fmt.Errorf("failed to create build execution in batch: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
