package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/test_suite/domain/models"
	"github.com/BennyEisner/test-results/internal/test_suite/domain/ports"
)

// SQLTestSuiteRepository implements the TestSuiteRepository interface
type SQLTestSuiteRepository struct {
	db *sql.DB
}

// NewSQLTestSuiteRepository creates a new SQL test suite repository
func NewSQLTestSuiteRepository(db *sql.DB) ports.TestSuiteRepository {
	return &SQLTestSuiteRepository{db: db}
}

// GetByID retrieves a test suite by its ID
func (r *SQLTestSuiteRepository) GetByID(ctx context.Context, id int64) (*models.TestSuite, error) {
	query := `SELECT id, project_id, name, parent_id, time FROM test_suites WHERE id = $1`

	var testSuite models.TestSuite
	var parentID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&testSuite.ID, &testSuite.ProjectID, &testSuite.Name, &parentID, &testSuite.Time,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get test suite by ID: %w", err)
	}

	if parentID.Valid {
		testSuite.ParentID = &parentID.Int64
	}

	return &testSuite, nil
}

// GetAllByProjectID retrieves all test suites for a project
func (r *SQLTestSuiteRepository) GetAllByProjectID(ctx context.Context, projectID int64) ([]*models.TestSuite, error) {
	query := `SELECT id, project_id, name, parent_id, time FROM test_suites WHERE project_id = $1 ORDER BY id`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test suites by project ID: %w", err)
	}
	defer rows.Close()

	var testSuites []*models.TestSuite
	for rows.Next() {
		var testSuite models.TestSuite
		var parentID sql.NullInt64

		err := rows.Scan(
			&testSuite.ID, &testSuite.ProjectID, &testSuite.Name, &parentID, &testSuite.Time,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan test suite: %w", err)
		}

		if parentID.Valid {
			testSuite.ParentID = &parentID.Int64
		}

		testSuites = append(testSuites, &testSuite)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating test suites: %w", err)
	}

	return testSuites, nil
}

// GetByName retrieves a test suite by its name within a project
func (r *SQLTestSuiteRepository) GetByName(ctx context.Context, projectID int64, name string) (*models.TestSuite, error) {
	query := `SELECT id, project_id, name, parent_id, time FROM test_suites WHERE project_id = $1 AND name = $2`

	var testSuite models.TestSuite
	var parentID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, projectID, name).Scan(
		&testSuite.ID, &testSuite.ProjectID, &testSuite.Name, &parentID, &testSuite.Time,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get test suite by name: %w", err)
	}

	if parentID.Valid {
		testSuite.ParentID = &parentID.Int64
	}

	return &testSuite, nil
}

// Create creates a new test suite
func (r *SQLTestSuiteRepository) Create(ctx context.Context, suite *models.TestSuite) error {
	query := `INSERT INTO test_suites (project_id, name, parent_id, time) VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		suite.ProjectID, suite.Name, suite.ParentID, suite.Time,
	).Scan(&suite.ID)
	if err != nil {
		return fmt.Errorf("failed to create test suite: %w", err)
	}

	return nil
}

// Update updates an existing test suite
func (r *SQLTestSuiteRepository) Update(ctx context.Context, id int64, name string) (*models.TestSuite, error) {
	query := `UPDATE test_suites SET name = $1 WHERE id = $2 RETURNING id, project_id, name, parent_id, time`

	var testSuite models.TestSuite
	var parentID sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, name, id).Scan(
		&testSuite.ID, &testSuite.ProjectID, &testSuite.Name, &parentID, &testSuite.Time,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update test suite: %w", err)
	}

	if parentID.Valid {
		testSuite.ParentID = &parentID.Int64
	}

	return &testSuite, nil
}

// Delete deletes a test suite by its ID
func (r *SQLTestSuiteRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM test_suites WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete test suite: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("test suite not found")
	}

	return nil
}
