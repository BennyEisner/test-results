package models

type TestSuite struct {
	ID        int64   `json:"id"`
	ProjectID int64   `json:"project_id"`
	Name      string  `json:"name"`
	ParentID  *int64  `json:"parent_id,omitempty"`
	Time      float64 `json:"time"`
}
