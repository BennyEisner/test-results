package models

import "time"

type Build struct {
	ID          int64     `json:"id"`
	TestSuiteID int64     `json:"test_suite_id"`
	BuildNumber string    `json:"build_number"`
	CIProvider  string    `json:"ci_provider"`
	CIURL       *string   `json:"ci_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
