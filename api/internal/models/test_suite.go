package models

type TestSuite struct {
	ID       int64   `json:"id"`
	BuildID  int64   `json:"build_id"`
	Name     string  `json:"name"`
	ParentID *int64  `json:"parent_id,omitempty"`
	Time     float64 `json:"time"`
}

