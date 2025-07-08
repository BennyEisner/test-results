package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/test_case/domain/models"
	"github.com/BennyEisner/test-results/internal/test_case/domain/ports"
)

// SQLTestCaseRepository implements the TestCaseRepository interface
type SQLTestCaseRepository struct {
	db *sql.DB
}

// NewSQLTestCaseRepository creates a new SQL test case repository
func NewSQLTestCaseRepository(db *sql.DB) ports.TestCaseRepository {
	return &SQLTestCaseRepository{db: db}
}

// GetByID retrieves a test case by its ID
func (r *SQLTestCaseRepository) GetByID(ctx context.Context, id int64) (*models.TestCase, error) {
	query := `SELECT id, suite_id, name, classname FROM test_cases WHERE id = $1`

	var testCase models.TestCase

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&testCase.ID, &testCase.SuiteID, &testCase.Name, &testCase.Classname,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get test case by ID: %w", err)
	}

	return &testCase, nil
}

// GetAllBySuiteID retrieves all test cases for a suite
func (r *SQLTestCaseRepository) GetAllBySuiteID(ctx context.Context, suiteID int64) ([]*models.TestCase, error) {
	query := `SELECT id, suite_id, name, classname FROM test_cases WHERE suite_id = $1 ORDER BY name`

	rows, err := r.db.QueryContext(ctx, query, suiteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test cases by suite ID: %w", err)
	}
	defer rows.Close()

	var testCases []*models.TestCase
	for rows.Next() {
		var testCase models.TestCase

		err := rows.Scan(
			&testCase.ID, &testCase.SuiteID, &testCase.Name, &testCase.Classname,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan test case: %w", err)
		}

		testCases = append(testCases, &testCase)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating test cases: %w", err)
	}

	return testCases, nil
}

// GetByName retrieves a test case by its name within a suite
func (r *SQLTestCaseRepository) GetByName(ctx context.Context, suiteID int64, name string) (*models.TestCase, error) {
	query := `SELECT id, suite_id, name, classname FROM test_cases WHERE suite_id = $1 AND name = $2`

	var testCase models.TestCase

	err := r.db.QueryRowContext(ctx, query, suiteID, name).Scan(
		&testCase.ID, &testCase.SuiteID, &testCase.Name, &testCase.Classname,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get test case by name: %w", err)
	}

	return &testCase, nil
}

// Create creates a new test case
func (r *SQLTestCaseRepository) Create(ctx context.Context, testCase *models.TestCase) error {
	query := `INSERT INTO test_cases (suite_id, name, classname) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		testCase.SuiteID, testCase.Name, testCase.Classname,
	).Scan(&testCase.ID)
	if err != nil {
		return fmt.Errorf("failed to create test case: %w", err)
	}

	return nil
}

// Update updates an existing test case
func (r *SQLTestCaseRepository) Update(ctx context.Context, id int64, name, classname string) (*models.TestCase, error) {
	query := `UPDATE test_cases SET name = $1, classname = $2 WHERE id = $3 RETURNING id, suite_id, name, classname`

	var testCase models.TestCase

	err := r.db.QueryRowContext(ctx, query, name, classname, id).Scan(
		&testCase.ID, &testCase.SuiteID, &testCase.Name, &testCase.Classname,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update test case: %w", err)
	}

	return &testCase, nil
}

// Delete deletes a test case by its ID
func (r *SQLTestCaseRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM test_cases WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete test case: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("test case not found")
	}

	return nil
}
