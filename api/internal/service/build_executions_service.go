package service

import (
	"database/sql"
	"fmt"
	"strings"

	// "github.com/BennyEisner/test-results/internal/handler" // No longer needed for DTOs
	"github.com/BennyEisner/test-results/internal/models"
	"github.com/BennyEisner/test-results/internal/utils"
)

// BuildExecutionServiceInterface defines the interface for build execution service operations.
type BuildExecutionServiceInterface interface {
	GetBuildExecutions(buildID int64) ([]models.BuildExecutionDetail, error)                                                     // Changed to models.BuildExecutionDetail
	CreateBuildExecutions(buildID int64, inputs []models.BuildExecutionInput) ([]models.BuildTestCaseExecution, []string, error) // Changed to models.BuildExecutionInput
	CheckBuildExists(buildID int64) (bool, error)
	CheckTestCaseExists(testCaseID int64) (bool, error)
}

// BuildExecutionService provides services related to build executions.
type BuildExecutionService struct {
	DB *sql.DB
}

// NewBuildExecutionService creates a new BuildExecutionService.
func NewBuildExecutionService(db *sql.DB) *BuildExecutionService {
	return &BuildExecutionService{DB: db}
}

// CheckBuildExists checks if a build with the given ID exists.
func (s *BuildExecutionService) CheckBuildExists(buildID int64) (bool, error) {
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM builds WHERE id = $1)", buildID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("database error checking build existence: %w", err)
	}
	return exists, nil
}

// CheckTestCaseExists checks if a test case with the given ID exists.
func (s *BuildExecutionService) CheckTestCaseExists(testCaseID int64) (bool, error) {
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM test_cases WHERE id = $1)", testCaseID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("database error checking test case existence: %w", err)
	}
	return exists, nil
}

// GetBuildExecutions fetches detailed execution results for a given build ID.
func (s *BuildExecutionService) GetBuildExecutions(buildID int64) ([]models.BuildExecutionDetail, error) { // Changed to models.BuildExecutionDetail
	query := `
		SELECT
			bte.id AS execution_id,
			bte.build_id,
			bte.test_case_id,
			tc.name AS test_case_name,
			tc.classname AS class_name,
			bte.status,
			bte.execution_time,
			bte.created_at,
			f.id AS failure_id,
			f.message AS failure_message,
			f.type AS failure_type,
			f.details AS failure_details
		FROM build_test_case_executions bte
		JOIN test_cases tc ON bte.test_case_id = tc.id
		LEFT JOIN failures f ON f.build_test_case_execution_id = bte.id
		WHERE bte.build_id = $1
		ORDER BY tc.name;
	`
	rows, err := s.DB.Query(query, buildID)
	if err != nil {
		return nil, fmt.Errorf("database error fetching build executions: %w", err)
	}
	defer rows.Close()

	var executions []models.BuildExecutionDetail // Changed to models.BuildExecutionDetail
	for rows.Next() {
		var detail models.BuildExecutionDetail // Changed to models.BuildExecutionDetail
		var failureID sql.NullInt64
		var failureMessage sql.NullString
		var failureType sql.NullString
		var failureDetails sql.NullString

		if err := rows.Scan(
			&detail.ExecutionID,
			&detail.BuildID,
			&detail.TestCaseID,
			&detail.TestCaseName,
			&detail.ClassName,
			&detail.Status,
			&detail.ExecutionTime,
			&detail.CreatedAt,
			&failureID,
			&failureMessage,
			&failureType,
			&failureDetails,
		); err != nil {
			return nil, fmt.Errorf("error scanning build execution: %w", err)
		}

		if failureID.Valid {
			detail.Failure = &models.Failure{
				ID:                       failureID.Int64,
				BuildTestCaseExecutionID: detail.ExecutionID,
				Message:                  utils.NullStringToStringPtr(failureMessage),
				Type:                     utils.NullStringToStringPtr(failureType),
				Details:                  utils.NullStringToStringPtr(failureDetails),
			}
		}
		executions = append(executions, detail)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating build execution rows: %w", err)
	}
	return executions, nil
}

// CreateBuildExecutions handles the creation of multiple build execution records.
// It returns a slice of successfully created executions, a slice of error messages for individual failures,
// and a general error if a non-recoverable issue occurs (e.g., cannot start transaction).
func (s *BuildExecutionService) CreateBuildExecutions(buildID int64, inputs []models.BuildExecutionInput) ([]models.BuildTestCaseExecution, []string, error) { // Changed to models.BuildExecutionInput
	var createdExecutions []models.BuildTestCaseExecution
	var processingErrors []string
	validStatuses := map[string]bool{"passed": true, "failed": true, "skipped": true, "error": true}

	for _, input := range inputs {
		tcExists, err := s.CheckTestCaseExists(input.TestCaseID)
		if err != nil {
			processingErrors = append(processingErrors, fmt.Sprintf("Error checking test case %d: %s", input.TestCaseID, err.Error()))
			continue
		}
		if !tcExists {
			processingErrors = append(processingErrors, fmt.Sprintf("Test case with ID %d not found.", input.TestCaseID))
			continue
		}

		statusToInsert := strings.ToLower(input.Status)
		if !validStatuses[statusToInsert] {
			processingErrors = append(processingErrors, fmt.Sprintf("Invalid status '%s' for test case ID %d. Must be one of: passed, failed, skipped, error.", input.Status, input.TestCaseID))
			continue
		}

		tx, err := s.DB.Begin()
		if err != nil {
			// This is a more global error, return it directly
			return createdExecutions, processingErrors, fmt.Errorf("failed to start database transaction: %w", err)
		}

		var executionID int64
		err = tx.QueryRow(
			`INSERT INTO build_test_case_executions (build_id, test_case_id, status, execution_time)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			buildID, input.TestCaseID, statusToInsert, input.ExecutionTime,
		).Scan(&executionID)

		if err != nil {
			tx.Rollback()
			processingErrors = append(processingErrors, fmt.Sprintf("Error inserting execution for test case %d: %s", input.TestCaseID, err.Error()))
			continue
		}

		if (statusToInsert == "failed" || statusToInsert == "error") &&
			(input.FailureMessage != nil || input.FailureType != nil || input.FailureDetails != nil) {
			_, err = tx.Exec(
				`INSERT INTO failures (build_test_case_execution_id, message, type, details)
				 VALUES ($1, $2, $3, $4)`,
				executionID, input.FailureMessage, input.FailureType, input.FailureDetails,
			)
			if err != nil {
				tx.Rollback()
				processingErrors = append(processingErrors, fmt.Sprintf("Error inserting failure details for test case %d (execution %d): %s", input.TestCaseID, executionID, err.Error()))
				continue
			}
		}

		if err := tx.Commit(); err != nil {
			processingErrors = append(processingErrors, fmt.Sprintf("Error committing transaction for test case %d: %s", input.TestCaseID, err.Error()))
			continue
		}

		// For simplicity, we're just adding a basic model here.
		// The handler might want to re-fetch or construct a more detailed DTO if needed.
		createdExecutions = append(createdExecutions, models.BuildTestCaseExecution{
			ID:            executionID,
			BuildID:       buildID,
			TestCaseID:    input.TestCaseID,
			Status:        statusToInsert,
			ExecutionTime: input.ExecutionTime,
			// CreatedAt would be set by DB, not available here without another query
		})
	}

	return createdExecutions, processingErrors, nil
}
