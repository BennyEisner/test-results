package models

type Failure struct {
	ID         int64   `json:"id"`
	TestCaseID int64   `json:"test_case_id"`
	Message    *string `json:"message,omitempty"`
	Type       *string `json:"type,omitempty"`
	Details    *string `json:"details,omitempty"`
}
