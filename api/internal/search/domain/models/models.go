package models

// SearchResult represents a search result
type SearchResult struct {
	Type      string `json:"type"`
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	ProjectID *int64 `json:"project_id,omitempty"`
	SuiteID   *int64 `json:"suite_id,omitempty"`
}
