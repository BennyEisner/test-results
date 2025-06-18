package models

import "time"

// BuildTestCaseExecution represents the execution of a specific test case within a specific build.
type BuildTestCaseExecution struct {
	ID            int64     `json:"id"`
	BuildID       int64     `json:"build_id"`
	TestCaseID    int64     `json:"test_case_id"`
	Status        string    `json:"status"`         // e.g., "passed", "failed", "skipped", "error"
	ExecutionTime float64   `json:"execution_time"` // Actual time taken for this specific execution
	CreatedAt     time.Time `json:"created_at"`
}

// BuildExecutionDetail is a DTO for returning detailed execution results.
// It's placed here to avoid import cycles between handler and service.
type BuildExecutionDetail struct {
	ExecutionID   int64     `json:"execution_id"` // ID from build_test_case_executions
	BuildID       int64     `json:"build_id"`
	TestCaseID    int64     `json:"test_case_id"`
	TestCaseName  string    `json:"test_case_name"`
	ClassName     string    `json:"class_name"`
	Status        string    `json:"status"`
	ExecutionTime float64   `json:"execution_time"`
	CreatedAt     time.Time `json:"created_at"`
	Failure       *Failure  `json:"failure,omitempty"` // Embed failure details if any
}

// BuildExecutionInput defines the structure for submitting a single test case execution result.
// It's placed here to avoid import cycles between handler and service.
type BuildExecutionInput struct {
	TestCaseID     int64   `json:"test_case_id"` // ID of the test case definition
	Status         string  `json:"status"`       // "passed", "failed", "skipped", "error"
	ExecutionTime  float64 `json:"execution_time"`
	FailureMessage *string `json:"failure_message,omitempty"`
	FailureType    *string `json:"failure_type,omitempty"`
	FailureDetails *string `json:"failure_details,omitempty"`
}
