package models

// SearchResult represents a search result
type SearchResult struct {
	Type string `json:"type"`
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}
