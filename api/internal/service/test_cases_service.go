package service

import (
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/models"
)

// TestCaseServiceInterface defines the interface for test case service operations.
type TestCaseServiceInterface interface {
	GetTestCasesBySuiteID(suiteID int64) ([]models.TestCase, error)
	CreateTestCase(suiteID int64, name string, classname string) (*models.TestCase, error)
	FindOrCreateTestCaseWithTx(tx *sql.Tx, suiteID int64, name string, classname string) (*models.TestCase, error) // New
	GetTestCaseByID(caseID int64) (*models.TestCase, error)
	CheckTestSuiteExists(suiteID int64) (bool, error) // To validate suite_id before creation
	GetMostFailedTests(projectID int64, limit int) ([]models.MostFailedTest, error)
}

type TestCaseService struct {
	DB *sql.DB
}

func NewTestCaseService(db *sql.DB) *TestCaseService {
	return &TestCaseService{DB: db}
}

// CheckTestSuiteExists checks if a test suite with the given ID exists.
// This is duplicated from TestSuiteService for now to avoid inter-service dependency at this stage.
// Ideally, this could be a shared utility or TestSuiteService could be injected.
func (s *TestCaseService) CheckTestSuiteExists(suiteID int64) (bool, error) {
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM test_suites WHERE id = $1)", suiteID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("database error checking test suite existence: %w", err)
	}
	return exists, nil
}

// GetTestCasesBySuiteID fetches all test cases for a given suiteID.
func (s *TestCaseService) GetTestCasesBySuiteID(suiteID int64) ([]models.TestCase, error) {
	rows, err := s.DB.Query("SELECT id, suite_id, name, classname FROM test_cases WHERE suite_id = $1 ORDER BY name", suiteID)
	if err != nil {
		return nil, fmt.Errorf("database error fetching test cases for suite %d: %w", suiteID, err)
	}
	defer rows.Close()

	cases := []models.TestCase{}
	for rows.Next() {
		var tc models.TestCase
		if err := rows.Scan(&tc.ID, &tc.SuiteID, &tc.Name, &tc.Classname); err != nil {
			return nil, fmt.Errorf("error scanning test case: %w", err)
		}
		cases = append(cases, tc)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating test case rows for suite %d: %w", suiteID, err)
	}
	return cases, nil
}

// CreateTestCase creates a new test case definition.
func (s *TestCaseService) CreateTestCase(suiteID int64, name string, classname string) (*models.TestCase, error) {
	var createdCase models.TestCase
	err := s.DB.QueryRow(
		"INSERT INTO test_cases(suite_id, name, classname) VALUES($1, $2, $3) RETURNING id, suite_id, name, classname",
		suiteID, name, classname,
	).Scan(&createdCase.ID, &createdCase.SuiteID, &createdCase.Name, &createdCase.Classname)

	if err != nil {
		// Consider checking for specific errors, e.g., foreign key violation for suite_id
		return nil, fmt.Errorf("database error creating test case: %w", err)
	}
	return &createdCase, nil
}

// GetTestCaseByID fetches a single test case by its ID.
func (s *TestCaseService) GetTestCaseByID(caseID int64) (*models.TestCase, error) {
	var tc models.TestCase
	err := s.DB.QueryRow("SELECT id, suite_id, name, classname FROM test_cases WHERE id = $1", caseID).Scan(
		&tc.ID, &tc.SuiteID, &tc.Name, &tc.Classname)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // Let handler decide on 404
		}
		return nil, fmt.Errorf("database error fetching test case by ID %d: %w", caseID, err)
	}
	return &tc, nil
}

// FindOrCreateTestCaseWithTx finds an existing test case by suite_id, name, and classname,
// or creates a new one if it doesn't exist, within an existing transaction.
func (s *TestCaseService) FindOrCreateTestCaseWithTx(tx *sql.Tx, suiteID int64, name string, classname string) (*models.TestCase, error) {
	var tc models.TestCase

	// Try to find existing test case
	// It's important that (suite_id, name, classname) is unique if we want this to be robust.
	// If classname can be empty or is not always reliable, the query might need adjustment or a unique constraint on (suite_id, name).
	// For now, assuming classname is significant.
	err := tx.QueryRow(
		"SELECT id, suite_id, name, classname FROM test_cases WHERE suite_id = $1 AND name = $2 AND classname = $3",
		suiteID, name, classname,
	).Scan(&tc.ID, &tc.SuiteID, &tc.Name, &tc.Classname)

	if err == nil {
		// Found existing test case
		return &tc, nil
	}

	if err != sql.ErrNoRows {
		// An actual error occurred during the query
		return nil, fmt.Errorf("error querying for test case (suite: %d, name: %s, class: %s) with tx: %w", suiteID, name, classname, err)
	}

	// Test case not found, so create it (err == sql.ErrNoRows)
	// Optional: Validate suiteID exists using CheckTestSuiteExistsWithTx if that method is added.
	// For now, relying on foreign key constraints or prior validation.

	var createdCase models.TestCase
	insertErr := tx.QueryRow(
		"INSERT INTO test_cases(suite_id, name, classname) VALUES($1, $2, $3) RETURNING id, suite_id, name, classname",
		suiteID, name, classname,
	).Scan(&createdCase.ID, &createdCase.SuiteID, &createdCase.Name, &createdCase.Classname)

	if insertErr != nil {
		return nil, fmt.Errorf("error creating test case (suite: %d, name: %s, class: %s) with tx: %w", suiteID, name, classname, insertErr)
	}

	return &createdCase, nil
}

func (s *TestCaseService) GetMostFailedTests(projectID int64, limit int) ([]models.MostFailedTest, error) {
	query := `
		SELECT
			tc.id AS test_case_id,
			tc.name,
			tc.classname,
			COUNT(tce.id) AS failure_count
		FROM
			test_cases tc
		JOIN
			build_test_case_executions tce ON tc.id = tce.test_case_id
		JOIN
			builds b ON tce.build_id = b.id
		JOIN
			test_suites ts ON b.test_suite_id = ts.id
		WHERE
			ts.project_id = $1 AND tce.status = 'failure'
		GROUP BY
			tc.id, tc.name, tc.classname
		ORDER BY
			failure_count DESC
		LIMIT $2;
	`

	rows, err := s.DB.Query(query, projectID, limit)
	if err != nil {
		return nil, fmt.Errorf("database error fetching most failed tests for project %d: %w", projectID, err)
	}
	defer rows.Close()

	var tests []models.MostFailedTest
	for rows.Next() {
		var test models.MostFailedTest
		if err := rows.Scan(&test.TestCaseID, &test.Name, &test.Classname, &test.FailureCount); err != nil {
			return nil, fmt.Errorf("error scanning most failed test: %w", err)
		}
		tests = append(tests, test)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating most failed test rows for project %d: %w", projectID, err)
	}

	return tests, nil
}
