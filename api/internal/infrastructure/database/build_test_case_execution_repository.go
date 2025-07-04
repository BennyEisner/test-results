package database

import (
	"context"
	"database/sql"

	"github.com/BennyEisner/test-results/internal/domain"
)

type SQLBuildTestCaseExecutionRepository struct {
	db *sql.DB
}

func NewSQLBuildTestCaseExecutionRepository(db *sql.DB) domain.BuildTestCaseExecutionRepository {
	return &SQLBuildTestCaseExecutionRepository{db: db}
}

func (r *SQLBuildTestCaseExecutionRepository) GetByID(ctx context.Context, id int64) (*domain.BuildTestCaseExecution, error) {
	query := `SELECT id, build_id, test_case_id, status, execution_time, created_at FROM build_test_case_executions WHERE id = $1`
	execution := &domain.BuildTestCaseExecution{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&execution.ID, &execution.BuildID, &execution.TestCaseID, &execution.Status, &execution.ExecutionTime, &execution.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return execution, nil
}

func (r *SQLBuildTestCaseExecutionRepository) GetAllByBuildID(ctx context.Context, buildID int64) ([]*domain.BuildExecutionDetail, error) {
	query := `
		SELECT 
			bte.id, bte.build_id, bte.test_case_id, tc.name, tc.classname, 
			bte.status, bte.execution_time, bte.created_at,
			f.message, f.type, f.details
		FROM build_test_case_executions bte
		JOIN test_cases tc ON bte.test_case_id = tc.id
		LEFT JOIN failures f ON f.build_test_case_execution_id = bte.id
		WHERE bte.build_id = $1
		ORDER BY bte.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, buildID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var executions []*domain.BuildExecutionDetail
	for rows.Next() {
		execution := &domain.BuildExecutionDetail{}
		var failureMessage, failureType, failureDetails sql.NullString
		if err := rows.Scan(&execution.ExecutionID, &execution.BuildID, &execution.TestCaseID, &execution.TestCaseName, &execution.ClassName, &execution.Status, &execution.ExecutionTime, &execution.CreatedAt, &failureMessage, &failureType, &failureDetails); err != nil {
			return nil, err
		}
		if failureMessage.Valid || failureType.Valid || failureDetails.Valid {
			execution.Failure = &domain.Failure{
				Message: &failureMessage.String,
				Type:    &failureType.String,
				Details: &failureDetails.String,
			}
		}
		executions = append(executions, execution)
	}
	return executions, nil
}

func (r *SQLBuildTestCaseExecutionRepository) Create(ctx context.Context, execution *domain.BuildTestCaseExecution) error {
	query := `INSERT INTO build_test_case_executions (build_id, test_case_id, status, execution_time, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, execution.BuildID, execution.TestCaseID, execution.Status, execution.ExecutionTime, execution.CreatedAt).Scan(&execution.ID)
}

func (r *SQLBuildTestCaseExecutionRepository) Update(ctx context.Context, id int64, execution *domain.BuildTestCaseExecution) (*domain.BuildTestCaseExecution, error) {
	query := `UPDATE build_test_case_executions SET build_id = $1, test_case_id = $2, status = $3, execution_time = $4 WHERE id = $5 RETURNING id, build_id, test_case_id, status, execution_time, created_at`
	updatedExecution := &domain.BuildTestCaseExecution{}
	if err := r.db.QueryRowContext(ctx, query, execution.BuildID, execution.TestCaseID, execution.Status, execution.ExecutionTime, id).Scan(&updatedExecution.ID, &updatedExecution.BuildID, &updatedExecution.TestCaseID, &updatedExecution.Status, &updatedExecution.ExecutionTime, &updatedExecution.CreatedAt); err != nil {
		return nil, err
	}
	return updatedExecution, nil
}

func (r *SQLBuildTestCaseExecutionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM build_test_case_executions WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
