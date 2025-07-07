package models

import "time"

// Failure represents a test failure
type Failure struct {
	ID          int64     `json:"id"`
	ExecutionID int64     `json:"execution_id"`
	Message     string    `json:"message"`
	Type        string    `json:"type"`
	Details     string    `json:"details"`
	CreatedAt   time.Time `json:"created_at"`
}
