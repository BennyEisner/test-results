package models

type TestCase struct {
	ID        int64  `json:"id"`
	SuiteID   int64  `json:"suite_id"` // The suite this test case definition belongs to
	Name      string `json:"name"`
	Classname string `json:"classname"`
	// Time (execution time) and Status are now part of BuildTestCaseExecution
}
