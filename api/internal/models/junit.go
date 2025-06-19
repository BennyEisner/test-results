package models

import "encoding/xml"

// JUnitFailure represents the <failure> element within a <testcase>.
// This element appears if a test assertion fails.
// Example: <testcase ...><failure message="expected true but was false" type="AssertionError">...</failure></testcase>
type JUnitFailure struct {
	XMLName xml.Name `xml:"failure"`                // Tells Go's XML parser that this struct maps to an XML element named "failure".
	Message string   `xml:"message,attr,omitempty"` // Maps to the 'message' attribute of <failure>. 'omitempty' means if the attribute is not present or empty, this field will be its zero value (empty string).
	Type    string   `xml:"type,attr,omitempty"`    // Maps to the 'type' attribute (e.g., "AssertionError").
	Value   string   `xml:",chardata"`              // Maps to the character data (the text content) inside the <failure>...</failure> tags. This often contains stack traces or more detailed failure info.
}

// JUnitError represents the <error> element within a <testcase>.
// This element appears if a test encounters an unexpected error during execution (e.g., a NullPointerException).
// Example: <testcase ...><error message="Unexpected null pointer" type="java.lang.NullPointerException">...</error></testcase>
type JUnitError struct {
	XMLName xml.Name `xml:"error"`                  // Maps to an XML element named "error".
	Message string   `xml:"message,attr,omitempty"` // Maps to the 'message' attribute.
	Type    string   `xml:"type,attr,omitempty"`    // Maps to the 'type' attribute.
	Value   string   `xml:",chardata"`              // Maps to the text content inside <error>...</error>.
}

// JUnitSkipped represents the <skipped> element within a <testcase>.
// This element appears if a test was intentionally skipped.
// Example: <testcase ...><skipped message="Test ignored" /></testcase>
// Note: Your <testsuite> sample also has a 'skipped' count attribute. This struct is for the <skipped> element *inside* a <testcase>.
type JUnitSkipped struct {
	XMLName xml.Name `xml:"skipped"`                // Maps to an XML element named "skipped".
	Message string   `xml:"message,attr,omitempty"` // Maps to the 'message' attribute, often explaining why it was skipped.
}

// JUnitTestCase represents a <testcase> element. Each <testcase> is an individual test.
// Example from your file: <testcase classname="tests.bdd.point.test_disputed_areas" name="test_atomic_v3_geo[28...]" time="0.545" />
type JUnitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`          // Maps to an XML element named "testcase".
	Classname string        `xml:"classname,attr"`    // Maps to the 'classname' attribute (e.g., "tests.bdd.point.test_disputed_areas").
	Name      string        `xml:"name,attr"`         // Maps to the 'name' attribute (e.g., "test_atomic_v3_geo[28...]").
	Time      float64       `xml:"time,attr"`         // Maps to the 'time' attribute (execution time in seconds). Go's XML parser can convert string attributes to float64.
	Failure   *JUnitFailure `xml:"failure,omitempty"` // Optional: A pointer to a JUnitFailure struct. If a <failure> sub-element exists, this will be populated. 'omitempty' means if there's no <failure> tag, this field will be nil.
	Error     *JUnitError   `xml:"error,omitempty"`   // Optional: A pointer to a JUnitError struct for <error> sub-elements.
	Skipped   *JUnitSkipped `xml:"skipped,omitempty"` // Optional: A pointer to a JUnitSkipped struct for <skipped> sub-elements.
}

// JUnitTestSuite represents a <testsuite> element. This groups multiple <testcase> elements.
// Example from your file: <testsuite name="pytest" errors="0" failures="0" skipped="8" tests="3009" time="116.202" timestamp="2025-06-18T06:10:38.315578+00:00" hostname="qatest-prod-useast1a-01">
type JUnitTestSuite struct {
	XMLName   xml.Name        `xml:"testsuite"`      // Maps to an XML element named "testsuite".
	Name      string          `xml:"name,attr"`      // Maps to the 'name' attribute (e.g., "pytest").
	Tests     int             `xml:"tests,attr"`     // Maps to the 'tests' attribute (total number of test cases in this suite).
	Failures  int             `xml:"failures,attr"`  // Maps to the 'failures' attribute.
	Errors    int             `xml:"errors,attr"`    // Maps to the 'errors' attribute.
	Skipped   int             `xml:"skipped,attr"`   // Maps to the 'skipped' attribute (count of skipped tests in this suite).
	Time      float64         `xml:"time,attr"`      // Maps to the 'time' attribute (total time for all tests in this suite).
	Timestamp string          `xml:"timestamp,attr"` // Maps to the 'timestamp' attribute (e.g., "2025-06-18T06:10:38.315578+00:00").
	Hostname  string          `xml:"hostname,attr"`  // Maps to the 'hostname' attribute (e.g., "qatest-prod-useast1a-01").
	TestCases []JUnitTestCase `xml:"testcase"`       // A slice to hold all <testcase> elements found within this <testsuite>.
}

// JUnitTestSuites represents the root <testsuites> element. This is the top-level element in your sample file.
// Example from your file: <testsuites name="pytest tests">
type JUnitTestSuites struct {
	XMLName    xml.Name         `xml:"testsuites"`          // Maps to an XML element named "testsuites".
	Name       string           `xml:"name,attr,omitempty"` // Maps to the 'name' attribute of <testsuites> (e.g., "pytest tests"). 'omitempty' for attributes means if it's not present in the XML, the field will be its zero value.
	TestSuites []JUnitTestSuite `xml:"testsuite"`           // A slice to hold all <testsuite> elements found within this <testsuites> root element.
	// Summary attributes like 'tests', 'failures' at the <testsuites> level are less common for the root <testsuites>
	// but can be added here if your files might contain them (e.g., Tests int `xml:"tests,attr,omitempty"`).
	// For now, this matches the structure from your sample where these summaries are per-<testsuite>.
}
