package models

type TestCase struct {
	ID       int64   `json:"id"`
	SuiteID  int64   `json:"suite_id"`
	Name     string  `json:"name"`
	Classname string `json:"classname"`
	Time     float64 `json:"time"`
	Status   string  `json:"status"` // "passed", "failed", "skipped"
}

