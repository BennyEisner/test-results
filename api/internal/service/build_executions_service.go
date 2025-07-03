package service

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/BennyEisner/test-results/internal/utils"
)

// BuildExecutionServiceInterface defines the interface for build execution service operations.
type BuildExecutionServiceInterface interface {
	GetBuildExecutions(buildID int64) ([]models.BuildExecutionDetail, error)
	CreateBuildExecutions(buildID int64, inputs []models.BuildExecutionInput) ([]models.BuildTestCaseExecution, []string, error)
	CreateBuildExecutionsWithTx(tx *sql.Tx, buildID int64, inputs []models.BuildExecutionInput) ([]models.BuildTestCaseExecution, []string, error) // New
	CheckBuildExists(buildID int64) (bool, error)
	CheckTestCaseExists(testCaseID int64) (bool, error)
	// Consider if CheckBuildExists and CheckTestCaseExists also need transactional versions
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

// validateBuildExecutionInput validates a single build execution input
func (s *BuildExecutionService) validateBuildExecutionInput(input models.BuildExecutionInput) error {
	validStatuses := map[string]bool{"passed": true, "failed": true, "skipped": true, "error": true}

	tcExists, err := s.CheckTestCaseExists(input.TestCaseID)
	if err != nil {
		return fmt.Errorf("error checking test case %d: %w", input.TestCaseID, err)
	}
	if !tcExists {
		return fmt.Errorf("test case with ID %d not found", input.TestCaseID)
	}

	statusToInsert := strings.ToLower(input.Status)
	if !validStatuses[statusToInsert] {
		return fmt.Errorf("invalid status '%s' for test case ID %d. Must be one of: passed, failed, skipped, error", input.Status, input.TestCaseID)
	}

	return nil
}

// shouldInsertFailureDetails determines if failure details should be inserted
func (s *BuildExecutionService) shouldInsertFailureDetails(status string, input models.BuildExecutionInput) bool {
	return (status == "failed" || status == "error") &&
		(input.FailureMessage != nil || input.FailureType != nil || input.FailureDetails != nil)
}

// createSingleExecution creates a single build execution with transaction handling
func (s *BuildExecutionService) createSingleExecution(buildID int64, input models.BuildExecutionInput) (*models.BuildTestCaseExecution, error) {
	statusToInsert := strings.ToLower(input.Status)

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to start database transaction: %w", err)
	}
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				// Log rollback error but don't return it since we're already returning an error
				_ = rbErr
			}
		}
	}()

	var executionID int64
	err = tx.QueryRow(
		`INSERT INTO build_test_case_executions (build_id, test_case_id, status, execution_time)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id`,
		buildID, input.TestCaseID, statusToInsert, input.ExecutionTime,
	).Scan(&executionID)

	if err != nil {
		return nil, fmt.Errorf("error inserting execution for test case %d: %w", input.TestCaseID, err)
	}

	// Handle failure details if needed
	if s.shouldInsertFailureDetails(statusToInsert, input) {
		if err := s.insertFailureDetails(tx, executionID, input); err != nil {
			return nil, fmt.Errorf("error inserting failure details for test case %d: %w", input.TestCaseID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction for test case %d: %w", input.TestCaseID, err)
	}

	return &models.BuildTestCaseExecution{
		ID:            executionID,
		BuildID:       buildID,
		TestCaseID:    input.TestCaseID,
		Status:        statusToInsert,
		ExecutionTime: input.ExecutionTime,
	}, nil
}

// insertFailureDetails inserts failure details into the database
func (s *BuildExecutionService) insertFailureDetails(tx *sql.Tx, executionID int64, input models.BuildExecutionInput) error {
	_, err := tx.Exec(
		`INSERT INTO failures (build_test_case_execution_id, message, type, details)
		 VALUES ($1, $2, $3, $4)`,
		executionID, input.FailureMessage, input.FailureType, input.FailureDetails,
	)
	return err
}

// CreateBuildExecutions handles the creation of multiple build execution records.
// It returns a slice of successfully created executions, a slice of error messages for individual failures,
// and a general error if a non-recoverable issue occurs (e.g., cannot start transaction).
func (s *BuildExecutionService) CreateBuildExecutions(buildID int64, inputs []models.BuildExecutionInput) ([]models.BuildTestCaseExecution, []string, error) {
	var createdExecutions []models.BuildTestCaseExecution
	var processingErrors []string

	for _, input := range inputs {
		// Validate input
		if err := s.validateBuildExecutionInput(input); err != nil {
			processingErrors = append(processingErrors, err.Error())
			continue
		}

		// Create execution
		execution, err := s.createSingleExecution(buildID, input)
		if err != nil {
			processingErrors = append(processingErrors, err.Error())
			continue
		}

		createdExecutions = append(createdExecutions, *execution)
	}

	return createdExecutions, processingErrors, nil
}

// CreateBuildExecutionsWithTx handles the creation of multiple build execution records within an existing transaction.
// It returns a slice of successfully created executions, a slice of error messages for individual failures,
// and a general error if a non-recoverable issue occurs during its own operations (but does not commit/rollback tx).
func (s *BuildExecutionService) CreateBuildExecutionsWithTx(tx *sql.Tx, buildID int64, inputs []models.BuildExecutionInput) ([]models.BuildTestCaseExecution, []string, error) {
	var createdExecutions []models.BuildTestCaseExecution
	var processingErrors []string
	validStatuses := map[string]bool{"passed": true, "failed": true, "skipped": true, "error": true}

	// Optional: Pre-check existence of all test cases in a batch if performance allows and it's critical.
	// For now, checking one by one.
	// Also, CheckTestCaseExists currently uses s.DB, not tx. If called here, it should ideally use tx.
	// For simplicity in this step, we'll assume test cases are validated before this call or rely on FK constraints.

	for _, input := range inputs {
		// It's better if the calling service (JUnitImportService) ensures test cases exist using its own transaction.
		// If we must check here:
		// var tcExists bool
		// err := tx.QueryRow("SELECT EXISTS(SELECT 1 FROM test_cases WHERE id = $1)", input.TestCaseID).Scan(&tcExists)
		// if err != nil {
		// 	processingErrors = append(processingErrors, fmt.Sprintf("Error checking test case %d with tx: %s", input.TestCaseID, err.Error()))
		// 	continue
		// }
		// if !tcExists {
		// 	processingErrors = append(processingErrors, fmt.Sprintf("Test case with ID %d not found (checked with tx).", input.TestCaseID))
		// 	continue
		// }

		statusToInsert := strings.ToLower(input.Status)
		if !validStatuses[statusToInsert] {
			processingErrors = append(processingErrors, fmt.Sprintf("Invalid status '%s' for test case ID %d. Must be one of: passed, failed, skipped, error.", input.Status, input.TestCaseID))
			continue
		}

		var executionID int64
		// Use the passed-in transaction 'tx'
		err := tx.QueryRow(
			`INSERT INTO build_test_case_executions (build_id, test_case_id, status, execution_time)
			 VALUES ($1, $2, $3, $4)
			 RETURNING id`,
			buildID, input.TestCaseID, statusToInsert, input.ExecutionTime,
		).Scan(&executionID)

		if err != nil {
			// Do not rollback here; let the caller manage the transaction.
			processingErrors = append(processingErrors, fmt.Sprintf("Error inserting execution for test case %d with tx: %s", input.TestCaseID, err.Error()))
			continue // Continue to try other executions in the batch
		}

		if (statusToInsert == "failed" || statusToInsert == "error") &&
			(input.FailureMessage != nil || input.FailureType != nil || input.FailureDetails != nil) {
			_, err = tx.Exec( // Use tx.Exec
				`INSERT INTO failures (build_test_case_execution_id, message, type, details)
				 VALUES ($1, $2, $3, $4)`,
				executionID, input.FailureMessage, input.FailureType, input.FailureDetails,
			)
			if err != nil {
				// Do not rollback here.
				processingErrors = append(processingErrors, fmt.Sprintf("Error inserting failure details for test case %d (execution %d) with tx: %s", input.TestCaseID, executionID, err.Error()))
				// This execution might be partially inserted (execution record created, failure not).
				// The caller needs to decide how to handle this based on the overall transaction strategy.
				continue
			}
		}

		// Do not commit here.

		createdExecutions = append(createdExecutions, models.BuildTestCaseExecution{
			ID:            executionID,
			BuildID:       buildID,
			TestCaseID:    input.TestCaseID,
			Status:        statusToInsert,
			ExecutionTime: input.ExecutionTime,
		})
	}

	// If any error occurred that should make the *entire batch operation* fail, return it as the third argument.
	// For now, we're collecting individual errors in processingErrors.
	return createdExecutions, processingErrors, nil
}
