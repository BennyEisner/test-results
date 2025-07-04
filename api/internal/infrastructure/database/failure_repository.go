package database

import (
	"context"
	"database/sql"

	"github.com/BennyEisner/test-results/internal/domain"
)

type SQLFailureRepository struct {
	db *sql.DB
}

func NewSQLFailureRepository(db *sql.DB) domain.FailureRepository {
	return &SQLFailureRepository{db: db}
}

func (r *SQLFailureRepository) GetByID(ctx context.Context, id int64) (*domain.Failure, error) {
	query := `SELECT id, build_test_case_execution_id, message, type, details FROM failures WHERE id = $1`
	failure := &domain.Failure{}
	var message, failureType, details sql.NullString
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&failure.ID, &failure.BuildTestCaseExecutionID, &message, &failureType, &details); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
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
	return failure, nil
}

func (r *SQLFailureRepository) GetByExecutionID(ctx context.Context, executionID int64) (*domain.Failure, error) {
	query := `SELECT id, build_test_case_execution_id, message, type, details FROM failures WHERE build_test_case_execution_id = $1`
	failure := &domain.Failure{}
	var message, failureType, details sql.NullString
	if err := r.db.QueryRowContext(ctx, query, executionID).Scan(&failure.ID, &failure.BuildTestCaseExecutionID, &message, &failureType, &details); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
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
	return failure, nil
}

func (r *SQLFailureRepository) Create(ctx context.Context, failure *domain.Failure) error {
	query := `INSERT INTO failures (build_test_case_execution_id, message, type, details) VALUES ($1, $2, $3, $4) RETURNING id`
	var message interface{} = nil
	var failureType interface{} = nil
	var details interface{} = nil
	if failure.Message != nil {
		message = *failure.Message
	}
	if failure.Type != nil {
		failureType = *failure.Type
	}
	if failure.Details != nil {
		details = *failure.Details
	}
	return r.db.QueryRowContext(ctx, query, failure.BuildTestCaseExecutionID, message, failureType, details).Scan(&failure.ID)
}

func (r *SQLFailureRepository) Update(ctx context.Context, id int64, failure *domain.Failure) (*domain.Failure, error) {
	query := `UPDATE failures SET build_test_case_execution_id = $1, message = $2, type = $3, details = $4 WHERE id = $5 RETURNING id, build_test_case_execution_id, message, type, details`
	var message interface{} = nil
	var failureType interface{} = nil
	var details interface{} = nil
	if failure.Message != nil {
		message = *failure.Message
	}
	if failure.Type != nil {
		failureType = *failure.Type
	}
	if failure.Details != nil {
		details = *failure.Details
	}
	updatedFailure := &domain.Failure{}
	var messageNull, failureTypeNull, detailsNull sql.NullString
	if err := r.db.QueryRowContext(ctx, query, failure.BuildTestCaseExecutionID, message, failureType, details, id).Scan(&updatedFailure.ID, &updatedFailure.BuildTestCaseExecutionID, &messageNull, &failureTypeNull, &detailsNull); err != nil {
		return nil, err
	}
	if messageNull.Valid {
		updatedFailure.Message = &messageNull.String
	}
	if failureTypeNull.Valid {
		updatedFailure.Type = &failureTypeNull.String
	}
	if detailsNull.Valid {
		updatedFailure.Details = &detailsNull.String
	}
	return updatedFailure, nil
}

func (r *SQLFailureRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM failures WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
