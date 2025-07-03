package service

import (
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/models"
)

// FailuresServiceInterface defines the interface for failures service operations.
type FailuresServiceInterface interface {
	GetFailuresByBuildID(buildID int64) ([]FailureWithTestCase, error)
}

// FailuresService handles failures related business logic
type FailuresService struct {
	db *sql.DB
}

// FailureWithTestCase represents a failure with associated test case informatio
type FailureWithTestCase struct {
	models.Failure
	TestCaseName      string   `json:"test_case_name"`
	TestCaseClassname string   `json:"test_case_classname"`
	ExecutionStatus   string   `json:"execution_status"`
	ExecutionTime     *float64 `json:"execution_time,omitempty"`
}

// NewFailuresService creates a new FailuresService.
func NewFailuresService(db *sql.DB) FailuresServiceInterface {
	return &FailuresService{db: db}
}

// GetFailuresByBuildID retrieves all failures for a specific build with test case details.
func (fs *FailuresService) GetFailuresByBuildID(buildID int64) ([]FailureWithTestCase, error) {
	query := `
		SELECT 
			f.id,
			f.build_test_case_execution_id,
			f.message,
			f.type,
			f.details,
			tc.name,
			tc.classname,
			btce.status,
			btce.execution_time
		FROM failures f
		JOIN build_test_case_executions btce ON f.build_test_case_execution_id = btce.id
		JOIN test_cases tc ON btce.test_case_id = tc.id
		WHERE btce.build_id = $1
		ORDER BY tc.classname, tc.name
	`

	rows, err := fs.db.Query(query, buildID)
	if err != nil {
		return nil, fmt.Errorf("failed to query failures for build %d: %w", buildID, err)
	}
	defer rows.Close()

	var failures []FailureWithTestCase
	for rows.Next() {
		var failure FailureWithTestCase
		err := rows.Scan(
			&failure.ID,
			&failure.BuildTestCaseExecutionID,
			&failure.Message,
			&failure.Type,
			&failure.Details,
			&failure.TestCaseName,
			&failure.TestCaseClassname,
			&failure.ExecutionStatus,
			&failure.ExecutionTime,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan failure row: %w", err)
		}
		failures = append(failures, failure)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating failure rows: %w", err)
	}

	return failures, nil
}
