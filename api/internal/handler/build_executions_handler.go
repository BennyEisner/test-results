package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/BennyEisner/test-results/internal/utils"
)

// BuildExecutionDetail is a DTO for returning detailed execution results.
type BuildExecutionDetail struct {
	ExecutionID   int64           `json:"execution_id"` // ID from build_test_case_executions
	BuildID       int64           `json:"build_id"`
	TestCaseID    int64           `json:"test_case_id"`
	TestCaseName  string          `json:"test_case_name"`
	ClassName     string          `json:"class_name"`
	Status        string          `json:"status"`
	ExecutionTime float64         `json:"execution_time"`
	CreatedAt     time.Time       `json:"created_at"`
	Failure       *models.Failure `json:"failure,omitempty"` // Embed failure details if any
}

// BuildExecutionInput defines the structure for submitting a single test case execution result.
type BuildExecutionInput struct {
	TestCaseID     int64   `json:"test_case_id"` // ID of the test case definition
	Status         string  `json:"status"`       // "passed", "failed", "skipped", "error"
	ExecutionTime  float64 `json:"execution_time"`
	FailureMessage *string `json:"failure_message,omitempty"`
	FailureType    *string `json:"failure_type,omitempty"`
	FailureDetails *string `json:"failure_details,omitempty"`
}

// HandleBuildExecutions handles GET and POST requests for test case executions of a specific build.
// Expected path: /api/builds/{buildId}/executions
func HandleBuildExecutions(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/builds/"), "/")

	// Expected: {buildId}/executions -> pathSegments = ["{buildId}", "executions"]
	if len(pathSegments) < 2 || pathSegments[0] == "" || pathSegments[1] != "executions" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL path for build executions. Expected /api/builds/{buildId}/executions")
		return
	}

	buildIDStr := pathSegments[0]
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID format: "+err.Error())
		return
	}

	// Check if the build exists
	var buildExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM builds WHERE id = $1)", buildID).Scan(&buildExists)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error checking build existence: "+err.Error())
		return
	}
	if !buildExists {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Build with ID %d not found", buildID))
		return
	}

	switch r.Method {
	case http.MethodGet:
		getBuildExecutions(w, r, buildID, db)
	case http.MethodPost:
		createBuildExecutions(w, r, buildID, db) // Placeholder for now
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed for this endpoint")
	}
}

func getBuildExecutions(w http.ResponseWriter, r *http.Request, buildID int64, db *sql.DB) {
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
	rows, err := db.Query(query, buildID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching build executions: "+err.Error())
		return
	}
	defer rows.Close()

	executions := []BuildExecutionDetail{}
	for rows.Next() {
		var detail BuildExecutionDetail
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
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning build execution: "+err.Error())
			return
		}

		if failureID.Valid { // If there was a failure for this execution
			detail.Failure = &models.Failure{
				ID:                       failureID.Int64,
				BuildTestCaseExecutionID: detail.ExecutionID, // This is the link
				Message:                  utils.NullStringToStringPtr(failureMessage),
				Type:                     utils.NullStringToStringPtr(failureType),
				Details:                  utils.NullStringToStringPtr(failureDetails),
			}
		}
		executions = append(executions, detail)
	}

	if err = rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating build execution rows: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, executions)
}

// createBuildExecutions handles POST requests to submit execution results for a build.
func createBuildExecutions(w http.ResponseWriter, r *http.Request, buildID int64, db *sql.DB) {
	var inputs []BuildExecutionInput
	if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	if len(inputs) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Request must contain at least one execution result.")
		return
	}

	// Validate status values
	validStatuses := map[string]bool{"passed": true, "failed": true, "skipped": true, "error": true}

	// Process each execution input

	var createdExecutions []models.BuildTestCaseExecution
	var processingErrors []string

	for _, input := range inputs {
		// Validate TestCaseID exists (optional, but good practice)
		var tcExists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM test_cases WHERE id = $1)", input.TestCaseID).Scan(&tcExists)
		if err != nil {
			processingErrors = append(processingErrors, fmt.Sprintf("Error checking test case %d: %s", input.TestCaseID, err.Error()))
			continue
		}
		if !tcExists {
			processingErrors = append(processingErrors, fmt.Sprintf("Test case with ID %d not found.", input.TestCaseID))
			continue
		}

		// Validate status
		statusToInsert := strings.ToLower(input.Status)
		if !validStatuses[statusToInsert] {
			processingErrors = append(processingErrors, fmt.Sprintf("Invalid status '%s' for test case ID %d. Must be one of: passed, failed, skipped, error.", input.Status, input.TestCaseID))
			continue
		}

		tx, err := db.Begin()
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Failed to start database transaction: "+err.Error())
			return // Abort all if transaction cannot be started
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

		// If status is 'failed' or 'error' and failure details are provided, insert into failures table
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
			// Rollback should have been called already if an exec error occurred,
			// but commit can also fail.
			processingErrors = append(processingErrors, fmt.Sprintf("Error committing transaction for test case %d: %s", input.TestCaseID, err.Error()))
			continue
		}

		// Successfully created, retrieve the full record for response (optional, but good for confirmation)
		createdExecutions = append(createdExecutions, models.BuildTestCaseExecution{ID: executionID, BuildID: buildID, TestCaseID: input.TestCaseID, Status: statusToInsert, ExecutionTime: input.ExecutionTime})
	}

	response := map[string]interface{}{
		"message":              fmt.Sprintf("%d execution results processed.", len(inputs)),
		"successful_creations": len(createdExecutions),
		"errors_encountered":   len(processingErrors),
	}
	if len(processingErrors) > 0 {
		response["error_details"] = processingErrors
		// Decide on overall status code if there are partial successes
		utils.RespondWithJSON(w, http.StatusMultiStatus, response)
	} else {
		utils.RespondWithJSON(w, http.StatusCreated, response)
	}
}
