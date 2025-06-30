package models

import "time"

type Build struct {
	ID            int64     `json:"id"`
	TestSuiteID   int64     `json:"test_suite_id"`
	ProjectID     int64     `json:"project_id"`
	BuildNumber   string    `json:"build_number"`
	CIProvider    string    `json:"ci_provider"`
	CIURL         *string   `json:"ci_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	TestCaseCount int64     `json:"test_case_count"`
}

type BuildDurationTrend struct {
	BuildNumber string    `json:"build_number"`
	Duration    float64   `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
}
