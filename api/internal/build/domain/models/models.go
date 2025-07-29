package models

import "time"

type Build struct {
	ID          int64     `json:"id"`
	ProjectID   int64     `json:"project_id"`
	SuiteID     int64     `json:"suite_id"`
	BuildNumber string    `json:"build_number"`
	Status      string    `json:"status"`
	Duration    float64   `json:"duration"`
	Timestamp   time.Time `json:"timestamp"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BuildDurationTrend struct {
	BuildNumber string    `json:"build_number"`
	Duration    float64   `json:"duration"`
	CreatedAt   time.Time `json:"created_at"`
}
