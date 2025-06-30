package models

type TestCase struct {
	ID        int64  `json:"id"`
	SuiteID   int64  `json:"suite_id"` // The suite this test case definition belongs to
	Name      string `json:"name"`
	Classname string `json:"classname"`
}

type MostFailedTest struct {
	TestCaseID   int64  `json:"test_case_id"`
	Name         string `json:"name"`
	Classname    string `json:"classname"`
	FailureCount int    `json:"failure_count"`
}
