package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// SQLFailureRepository implements the FailureRepository interface
type SQLFailureRepository struct {
	db *sql.DB
}

// NewSQLFailureRepository creates a new SQL failure repository
func NewSQLFailureRepository(db *sql.DB) ports.FailureRepository {
	return &SQLFailureRepository{db: db}
}

// GetByID retrieves a failure by its ID
func (r *SQLFailureRepository) GetByID(ctx context.Context, id int64) (*models.Failure, error) {
	query := `SELECT id, build_test_case_execution_id, message, type, details FROM failures WHERE id = $1`

	var failure models.Failure
	var message, failureType, details sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&failure.ID, &failure.BuildTestCaseExecutionID, &message, &failureType, &details,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get failure by ID: %w", err)
	}

	if message.Valid {
		failure.Message = &message.String
	}
	if failureType.Valid {
		failure.Type = &failureType.String
	}
	if details.Valid {
		failure.Details = &details.String
	}

	return &failure, nil
}

// GetByExecutionID retrieves a failure by execution ID
func (r *SQLFailureRepository) GetByExecutionID(ctx context.Context, executionID int64) (*models.Failure, error) {
	query := `SELECT id, build_test_case_execution_id, message, type, details FROM failures WHERE build_test_case_execution_id = $1`

	var failure models.Failure
	var message, failureType, details sql.NullString

	err := r.db.QueryRowContext(ctx, query, executionID).Scan(
		&failure.ID, &failure.BuildTestCaseExecutionID, &message, &failureType, &details,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get failure by execution ID: %w", err)
	}

	if message.Valid {
		failure.Message = &message.String
	}
	if failureType.Valid {
		failure.Type = &failureType.String
	}
	if details.Valid {
		failure.Details = &details.String
	}

	return &failure, nil
}

// Create creates a new failure
func (r *SQLFailureRepository) Create(ctx context.Context, failure *models.Failure) error {
	query := `INSERT INTO failures (build_test_case_execution_id, message, type, details) VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		failure.BuildTestCaseExecutionID, failure.Message, failure.Type, failure.Details,
	).Scan(&failure.ID)
	if err != nil {
		return fmt.Errorf("failed to create failure: %w", err)
	}

	return nil
}

// Update updates an existing failure
func (r *SQLFailureRepository) Update(ctx context.Context, id int64, failure *models.Failure) (*models.Failure, error) {
	query := `UPDATE failures SET message = $1, type = $2, details = $3 WHERE id = $4 RETURNING id, build_test_case_execution_id, message, type, details`

	var updatedFailure models.Failure
	var message, failureType, details sql.NullString

	err := r.db.QueryRowContext(ctx, query, failure.Message, failure.Type, failure.Details, id).Scan(
		&updatedFailure.ID, &updatedFailure.BuildTestCaseExecutionID, &message, &failureType, &details,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update failure: %w", err)
	}

	if message.Valid {
		updatedFailure.Message = &message.String
	}
	if failureType.Valid {
		updatedFailure.Type = &failureType.String
	}
	if details.Valid {
		updatedFailure.Details = &details.String
	}

	return &updatedFailure, nil
}

// Delete deletes a failure by its ID
func (r *SQLFailureRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM failures WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete failure: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("failure not found")
	}

	return nil
}
