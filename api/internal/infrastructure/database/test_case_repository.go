package database

import (
	"context"
	"database/sql"

	"github.com/BennyEisner/test-results/internal/domain"
)

type SQLTestCaseRepository struct {
	db *sql.DB
}

func NewSQLTestCaseRepository(db *sql.DB) domain.TestCaseRepository {
	return &SQLTestCaseRepository{db: db}
}

func (r *SQLTestCaseRepository) GetByID(ctx context.Context, id int64) (*domain.TestCase, error) {
	query := `SELECT id, suite_id, name, classname FROM test_cases WHERE id = $1`
	testCase := &domain.TestCase{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&testCase.ID, &testCase.SuiteID, &testCase.Name, &testCase.Classname); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return testCase, nil
}

func (r *SQLTestCaseRepository) GetAllBySuiteID(ctx context.Context, suiteID int64) ([]*domain.TestCase, error) {
	query := `SELECT id, suite_id, name, classname FROM test_cases WHERE suite_id = $1 ORDER BY name`
	rows, err := r.db.QueryContext(ctx, query, suiteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var testCases []*domain.TestCase
	for rows.Next() {
		testCase := &domain.TestCase{}
		if err := rows.Scan(&testCase.ID, &testCase.SuiteID, &testCase.Name, &testCase.Classname); err != nil {
			return nil, err
		}
		testCases = append(testCases, testCase)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return testCases, nil
}

func (r *SQLTestCaseRepository) GetByName(ctx context.Context, suiteID int64, name string) (*domain.TestCase, error) {
	query := `SELECT id, suite_id, name, classname FROM test_cases WHERE suite_id = $1 AND name = $2`
	testCase := &domain.TestCase{}
	if err := r.db.QueryRowContext(ctx, query, suiteID, name).Scan(&testCase.ID, &testCase.SuiteID, &testCase.Name, &testCase.Classname); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return testCase, nil
}

func (r *SQLTestCaseRepository) Create(ctx context.Context, tc *domain.TestCase) error {
	query := `INSERT INTO test_cases (suite_id, name, classname) VALUES ($1, $2, $3) RETURNING id`
	return r.db.QueryRowContext(ctx, query, tc.SuiteID, tc.Name, tc.Classname).Scan(&tc.ID)
}

func (r *SQLTestCaseRepository) Update(ctx context.Context, id int64, name, classname string) (*domain.TestCase, error) {
	query := `UPDATE test_cases SET name = $1, classname = $2 WHERE id = $3 RETURNING id, suite_id, name, classname`
	testCase := &domain.TestCase{}
	if err := r.db.QueryRowContext(ctx, query, name, classname, id).Scan(&testCase.ID, &testCase.SuiteID, &testCase.Name, &testCase.Classname); err != nil {
		return nil, err
	}
	return testCase, nil
}

func (r *SQLTestCaseRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM test_cases WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
