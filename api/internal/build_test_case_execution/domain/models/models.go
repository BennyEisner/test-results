package models

import "time"

// BuildTestCaseExecution represents a test case execution within a build
type BuildTestCaseExecution struct {
	ID            int64     `json:"id"`
	BuildID       int64     `json:"build_id"`
	TestCaseID    int64     `json:"test_case_id"`
	Status        string    `json:"status"`
	ExecutionTime float64   `json:"execution_time"`
	CreatedAt     time.Time `json:"created_at"`
}

// BuildExecution represents a test case execution within a build
type BuildExecution struct {
	ID            int64     `json:"id"`
	BuildID       int64     `json:"build_id"`
	TestCaseID    int64     `json:"test_case_id"`
	Status        string    `json:"status"`
	ExecutionTime float64   `json:"execution_time"`
	CreatedAt     time.Time `json:"created_at"`
}

// BuildExecutionDetail represents detailed build execution information
type BuildExecutionDetail struct {
	ExecutionID   int64     `json:"execution_id"`
	BuildID       int64     `json:"build_id"`
	TestCaseID    int64     `json:"test_case_id"`
	TestCaseName  string    `json:"test_case_name"`
	ClassName     string    `json:"class_name"`
	Status        string    `json:"status"`
	ExecutionTime float64   `json:"execution_time"`
	CreatedAt     time.Time `json:"created_at"`
	Failure       *Failure  `json:"failure,omitempty"`
}

// BuildExecutionInput represents input for creating a build execution
type BuildExecutionInput struct {
	TestCaseID     int64   `json:"test_case_id"`
	Status         string  `json:"status"`
	ExecutionTime  float64 `json:"execution_time"`
	FailureMessage *string `json:"failure_message,omitempty"`
	FailureType    *string `json:"failure_type,omitempty"`
	FailureDetails *string `json:"failure_details,omitempty"`
}

// Failure represents a test failure (minimal definition for build_test_case_execution domain)
type Failure struct {
	Message string `json:"message,omitempty"`
	Type    string `json:"type,omitempty"`
	Details string `json:"details,omitempty"`
}
