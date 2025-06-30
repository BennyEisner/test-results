package utils

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Build struct {
	ID            int       `json:"id"`
	TestSuiteID   int       `json:"test_suite_id"`
	ProjectID     int       `json:"project_id"`
	BuildNumber   string    `json:"build_number"`
	CIProvider    string    `json:"ci_provider"`
	CIURL         string    `json:"ci_url,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	TestCaseCount int       `json:"test_case_count,omitempty"`
	Duration      *float64  `json:"duration,omitempty"`
}

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, map[string]string{"error": message})
}

// NullStringToStringPtr converts a sql.NullString to a *string.
// If the NullString is valid, it returns a pointer to its string value.
// Otherwise, it returns nil.
func NullStringToStringPtr(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

// NullInt64ToIntPtr converts a sql.NullInt64 to a *int64.
// If the NullInt64 is valid, it returns a pointer to its int64 value.
// Otherwise, it returns nil.
func NullInt64ToIntPtr(ni sql.NullInt64) *int64 {
	if ni.Valid {
		return &ni.Int64
	}
	return nil
}

// NullFloat64ToFloat64Ptr converts a sql.NullFloat64 to a *float64.
// If the NullFloat64 is valid, it returns a pointer to its float64 value.
// Otherwise, it returns nil.
func NullFloat64ToFloat64Ptr(nf sql.NullFloat64) *float64 {
	if nf.Valid {
		return &nf.Float64
	}
	return nil
}
