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
