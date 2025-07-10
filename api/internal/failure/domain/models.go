package models

// Failure represents a test failure
type Failure struct {
	ID                       int64   `json:"id"`
	BuildTestCaseExecutionID int64   `json:"build_test_case_execution_id"`
	Message                  *string `json:"message,omitempty"`
	Type                     *string `json:"type,omitempty"`
	Details                  *string `json:"details,omitempty"`
}
