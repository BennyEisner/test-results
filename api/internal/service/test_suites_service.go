package service

import (
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/models"
)

// TestSuiteServiceInterface defines the interface for test suite service operations.
type TestSuiteServiceInterface interface {
	GetTestSuitesByProjectID(projectID int64) ([]models.TestSuite, error)
	CreateTestSuite(projectID int64, name string, parentID *int64, time float64) (*models.TestSuite, error)
	GetTestSuiteByID(suiteID int64) (*models.TestSuite, error)
	GetProjectTestSuiteByID(projectID int64, suiteID int64) (*models.TestSuite, error)
	CheckProjectExists(projectID int64) (bool, error)
	CheckTestSuiteExists(suiteID int64) (bool, error) // For validating parent_id
}

// TestSuiteService provides services related to test suites.
type TestSuiteService struct {
	DB *sql.DB
}

// NewTestSuiteService creates a new TestSuiteService.
func NewTestSuiteService(db *sql.DB) *TestSuiteService {
	return &TestSuiteService{DB: db}
}

// CheckProjectExists checks if a project with the given ID exists.
func (s *TestSuiteService) CheckProjectExists(projectID int64) (bool, error) {
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1)", projectID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("database error checking project existence: %w", err)
	}
	return exists, nil
}

// CheckTestSuiteExists checks if a test suite with the given ID exists.
func (s *TestSuiteService) CheckTestSuiteExists(suiteID int64) (bool, error) {
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM test_suites WHERE id = $1)", suiteID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("database error checking test suite existence: %w", err)
	}
	return exists, nil
}

// GetTestSuitesByProjectID fetches all test suites for a given projectID.
func (s *TestSuiteService) GetTestSuitesByProjectID(projectID int64) ([]models.TestSuite, error) {
	rows, err := s.DB.Query("SELECT id, project_id, name, parent_id, time FROM test_suites WHERE project_id = $1 ORDER BY name", projectID)
	if err != nil {
		return nil, fmt.Errorf("database error fetching test suites for project %d: %w", projectID, err)
	}
	defer rows.Close()

	suites := []models.TestSuite{}
	for rows.Next() {
		var ts models.TestSuite
		var parentID sql.NullInt64
		if err := rows.Scan(&ts.ID, &ts.ProjectID, &ts.Name, &parentID, &ts.Time); err != nil {
			return nil, fmt.Errorf("error scanning test suite: %w", err)
		}
		if parentID.Valid {
			ts.ParentID = &parentID.Int64
		}
		suites = append(suites, ts)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating test suite rows for project %d: %w", projectID, err)
	}
	return suites, nil
}

// CreateTestSuite creates a new test suite.
func (s *TestSuiteService) CreateTestSuite(projectID int64, name string, parentIDIn *int64, time float64) (*models.TestSuite, error) {
	var createdSuite models.TestSuite
	var parentIDArg sql.NullInt64

	if parentIDIn != nil {
		// Optional: Validate parent suite exists and belongs to the same project
		// For now, assuming frontend/caller ensures this or DB foreign key handles it.
		parentIDArg = sql.NullInt64{Int64: *parentIDIn, Valid: true}
	}

	err := s.DB.QueryRow(
		"INSERT INTO test_suites(project_id, name, parent_id, time) VALUES($1, $2, $3, $4) RETURNING id, project_id, name, parent_id, time",
		projectID, name, parentIDArg, time,
	).Scan(&createdSuite.ID, &createdSuite.ProjectID, &createdSuite.Name, &parentIDArg, &createdSuite.Time)

	if err != nil {
		return nil, fmt.Errorf("database error creating test suite: %w", err)
	}
	if parentIDArg.Valid {
		createdSuite.ParentID = &parentIDArg.Int64
	}

	return &createdSuite, nil
}

// GetTestSuiteByID fetches a single test suite by its ID.
func (s *TestSuiteService) GetTestSuiteByID(suiteID int64) (*models.TestSuite, error) {
	var ts models.TestSuite
	var parentID sql.NullInt64
	err := s.DB.QueryRow("SELECT id, project_id, name, parent_id, time FROM test_suites WHERE id = $1", suiteID).Scan(
		&ts.ID, &ts.ProjectID, &ts.Name, &parentID, &ts.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // Let handler decide on 404
		}
		return nil, fmt.Errorf("database error fetching test suite by ID %d: %w", suiteID, err)
	}
	if parentID.Valid {
		ts.ParentID = &parentID.Int64
	}
	return &ts, nil
}

// GetProjectTestSuiteByID fetches a specific test suite by its ID and projectID.
func (s *TestSuiteService) GetProjectTestSuiteByID(projectID int64, suiteID int64) (*models.TestSuite, error) {
	var ts models.TestSuite
	var parentID sql.NullInt64

	// Service can rely on its own CheckProjectExists or let the query handle it.
	// For now, the query implicitly checks project membership.
	err := s.DB.QueryRow("SELECT id, project_id, name, parent_id, time FROM test_suites WHERE id = $1 AND project_id = $2", suiteID, projectID).Scan(
		&ts.ID, &ts.ProjectID, &ts.Name, &parentID, &ts.Time)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // Let handler decide on 404
		}
		return nil, fmt.Errorf("database error fetching test suite ID %d for project ID %d: %w", suiteID, projectID, err)
	}

	if parentID.Valid {
		ts.ParentID = &parentID.Int64
	}
	return &ts, nil
}
