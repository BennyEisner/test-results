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
	GetTestCaseByID(caseID int64) (*models.TestCase, error)
	CheckTestSuiteExists(suiteID int64) (bool, error) // To validate suite_id before creation
}

// TestCaseService provides services related to test cases.
type TestCaseService struct {
	DB *sql.DB
	// testSuiteService TestSuiteServiceInterface // Could be injected if we want to use its CheckTestSuiteExists
}

// NewTestCaseService creates a new TestCaseService.
// func NewTestCaseService(db *sql.DB, tsService TestSuiteServiceInterface) *TestCaseService {
// return &TestCaseService{DB: db, testSuiteService: tsService}
// }
// Simplified version for now:
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
