package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

type Project struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Build struct {
	ID          int       `json:"id"`
	TestSuiteID int       `json:"test_suite_id"`
	BuildNumber string    `json:"build_number"`
	CIProvider  string    `json:"ci_provider"`
	CIURL       string    `json:"ci_url, omitempty"`
	CreatedAt   time.Time `json:"created_at, omitempty"`
}

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func RespondWithError(w http.ResponseWriter, status int, message string) {
	RespondWithJSON(w, status, map[string]string{"error": message})
}
