package models

// JUnitTestSuites represents JUnit XML data
type JUnitTestSuites struct {
	Name       string           `xml:"name,attr"`
	TestSuites []JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite represents a test suite in JUnit XML
type JUnitTestSuite struct {
	Name      string          `xml:"name,attr"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

// JUnitTestCase represents a test case in JUnit XML
type JUnitTestCase struct {
	Name      string        `xml:"name,attr"`
	Classname string        `xml:"classname,attr"`
	Time      float64       `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
	Error     *JUnitError   `xml:"error,omitempty"`
	Skipped   *JUnitSkipped `xml:"skipped,omitempty"`
}

// JUnitFailure represents a test failure in JUnit XML
type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Value   string `xml:",chardata"`
}

// JUnitError represents a test error in JUnit XML
type JUnitError struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Value   string `xml:",chardata"`
}

// JUnitSkipped represents a skipped test in JUnit XML
type JUnitSkipped struct {
	Message string `xml:"message,attr"`
}
