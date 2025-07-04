package models

// TestCase represents a test case within a test suite
type TestCase struct {
	ID        int64  `json:"id"`
	SuiteID   int64  `json:"suite_id"`
	Name      string `json:"name"`
	Classname string `json:"classname"`
}
