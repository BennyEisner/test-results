package domain

import (
	"context"
	"errors"
	"time"
)

// =============================================================================
// DOMAIN MODELS (Core Business Entities)
// =============================================================================

// Project represents a core business entity
type Project struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// TestSuite represents a test suite within a project
type TestSuite struct {
	ID        int64   `json:"id"`
	ProjectID int64   `json:"project_id"`
	Name      string  `json:"name"`
	ParentID  *int64  `json:"parent_id,omitempty"`
	Time      float64 `json:"time"`
}

// TestCase represents a test case within a test suite
type TestCase struct {
	ID        int64  `json:"id"`
	SuiteID   int64  `json:"suite_id"`
	Name      string `json:"name"`
	Classname string `json:"classname"`
}

// Build represents a build execution
type Build struct {
	ID            int64      `json:"id"`
	TestSuiteID   int64      `json:"test_suite_id"`
	ProjectID     int64      `json:"project_id"`
	BuildNumber   string     `json:"build_number"`
	CIProvider    string     `json:"ci_provider"`
	CIURL         *string    `json:"ci_url,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	StartedAt     *time.Time `json:"started_at,omitempty"`
	EndedAt       *time.Time `json:"ended_at,omitempty"`
	Duration      *float64   `json:"duration,omitempty"`
	TestCaseCount int64      `json:"test_case_count"`
}

// BuildExecution represents a test case execution within a build
type BuildExecution struct {
	ID            int64     `json:"id"`
	BuildID       int64     `json:"build_id"`
	TestCaseID    int64     `json:"test_case_id"`
	Status        string    `json:"status"`
	ExecutionTime float64   `json:"execution_time"`
	CreatedAt     time.Time `json:"created_at"`
}

// Failure represents a test failure
type Failure struct {
	ID                       int64   `json:"id"`
	BuildTestCaseExecutionID int64   `json:"build_test_case_execution_id"`
	Message                  *string `json:"message,omitempty"`
	Type                     *string `json:"type,omitempty"`
	Details                  *string `json:"details,omitempty"`
}

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

// =============================================================================
// BUILD TEST CASE EXECUTION DOMAIN
// =============================================================================

type BuildTestCaseExecution struct {
	ID            int64     `json:"id"`
	BuildID       int64     `json:"build_id"`
	TestCaseID    int64     `json:"test_case_id"`
	Status        string    `json:"status"`
	ExecutionTime float64   `json:"execution_time"`
	CreatedAt     time.Time `json:"created_at"`
}

type BuildExecutionDetail struct {
	ExecutionID   int64     `json:"execution_id"`
	BuildID       int64     `json:"build_id"`
	TestCaseID    int64     `json:"test_case_id"`
	TestCaseName  string    `json:"test_case_name"`
	ClassName     string    `json:"class_name"`
	Status        string    `json:"status"`
	ExecutionTime float64   `json:"execution_time"`
	CreatedAt     time.Time `json:"created_at"`
	Failure       *Failure  `json:"failure,omitempty"`
}

type BuildExecutionInput struct {
	TestCaseID     int64   `json:"test_case_id"`
	Status         string  `json:"status"`
	ExecutionTime  float64 `json:"execution_time"`
	FailureMessage *string `json:"failure_message,omitempty"`
	FailureType    *string `json:"failure_type,omitempty"`
	FailureDetails *string `json:"failure_details,omitempty"`
}

type BuildTestCaseExecutionRepository interface {
	GetByID(ctx context.Context, id int64) (*BuildTestCaseExecution, error)
	GetAllByBuildID(ctx context.Context, buildID int64) ([]*BuildExecutionDetail, error)
	Create(ctx context.Context, execution *BuildTestCaseExecution) error
	Update(ctx context.Context, id int64, execution *BuildTestCaseExecution) (*BuildTestCaseExecution, error)
	Delete(ctx context.Context, id int64) error
}

type BuildTestCaseExecutionService interface {
	GetExecutionByID(ctx context.Context, id int64) (*BuildTestCaseExecution, error)
	GetExecutionsByBuildID(ctx context.Context, buildID int64) ([]*BuildExecutionDetail, error)
	CreateExecution(ctx context.Context, buildID int64, input *BuildExecutionInput) (*BuildTestCaseExecution, error)
	UpdateExecution(ctx context.Context, id int64, execution *BuildTestCaseExecution) (*BuildTestCaseExecution, error)
	DeleteExecution(ctx context.Context, id int64) error
}

// =============================================================================
// INPUT PORTS (Primary/Driving Adapters)
// =============================================================================

// ProjectService defines the business logic for project operations
type ProjectService interface {
	GetProjectByID(ctx context.Context, id int64) (*Project, error)
	GetAllProjects(ctx context.Context) ([]*Project, error)
	CreateProject(ctx context.Context, name string) (*Project, error)
	UpdateProject(ctx context.Context, id int64, name string) (*Project, error)
	DeleteProject(ctx context.Context, id int64) error
	GetProjectByName(ctx context.Context, name string) (*Project, error)
}

// BuildService defines the business logic for build operations
type BuildService interface {
	GetBuildByID(ctx context.Context, id int64) (*Build, error)
	GetBuildsByProjectID(ctx context.Context, projectID int64) ([]*Build, error)
	GetBuildsByTestSuiteID(ctx context.Context, suiteID int64) ([]*Build, error)
	CreateBuild(ctx context.Context, build *Build) (*Build, error)
	UpdateBuild(ctx context.Context, id int64, build *Build) (*Build, error)
	DeleteBuild(ctx context.Context, id int64) error
}

// TestSuiteService defines the business logic for test suite operations
type TestSuiteService interface {
	GetTestSuiteByID(ctx context.Context, id int64) (*TestSuite, error)
	GetTestSuitesByProjectID(ctx context.Context, projectID int64) ([]*TestSuite, error)
	GetTestSuiteByName(ctx context.Context, projectID int64, name string) (*TestSuite, error)
	CreateTestSuite(ctx context.Context, projectID int64, name string, parentID *int64) (*TestSuite, error)
	UpdateTestSuite(ctx context.Context, id int64, name string) (*TestSuite, error)
	DeleteTestSuite(ctx context.Context, id int64) error
}

// TestCaseService defines the business logic for test case operations
type TestCaseService interface {
	GetTestCaseByID(ctx context.Context, id int64) (*TestCase, error)
	GetTestCasesBySuiteID(ctx context.Context, suiteID int64) ([]*TestCase, error)
	GetTestCaseByName(ctx context.Context, suiteID int64, name string) (*TestCase, error)
	CreateTestCase(ctx context.Context, suiteID int64, name, classname string) (*TestCase, error)
	UpdateTestCase(ctx context.Context, id int64, name, classname string) (*TestCase, error)
	DeleteTestCase(ctx context.Context, id int64) error
}

// BuildExecutionService defines the business logic for build execution operations
type BuildExecutionService interface {
	GetBuildExecutions(ctx context.Context, buildID int64) ([]*BuildExecution, error)
	CreateBuildExecutions(ctx context.Context, buildID int64, executions []*BuildExecution) error
}

// JUnitImportService defines the business logic for JUnit import operations
type JUnitImportService interface {
	ProcessJUnitData(ctx context.Context, projectID int64, suiteID int64, junitData *JUnitTestSuites) (*Build, error)
}

// =============================================================================
// OUTPUT PORTS (Secondary/Driven Adapters)
// =============================================================================

// ProjectRepository defines the data access contract for projects
type ProjectRepository interface {
	GetByID(ctx context.Context, id int64) (*Project, error)
	GetAll(ctx context.Context) ([]*Project, error)
	GetByName(ctx context.Context, name string) (*Project, error)
	Create(ctx context.Context, p *Project) error
	Update(ctx context.Context, id int64, name string) (*Project, error)
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context) (int, error)
}

// BuildRepository defines the data access contract for builds
type BuildRepository interface {
	GetByID(ctx context.Context, id int64) (*Build, error)
	GetAllByProjectID(ctx context.Context, projectID int64) ([]*Build, error)
	GetAllByTestSuiteID(ctx context.Context, suiteID int64) ([]*Build, error)
	Create(ctx context.Context, build *Build) error
	Update(ctx context.Context, id int64, build *Build) (*Build, error)
	Delete(ctx context.Context, id int64) error
}

// TestSuiteRepository defines the data access contract for test suites
type TestSuiteRepository interface {
	GetByID(ctx context.Context, id int64) (*TestSuite, error)
	GetAllByProjectID(ctx context.Context, projectID int64) ([]*TestSuite, error)
	GetByName(ctx context.Context, projectID int64, name string) (*TestSuite, error)
	Create(ctx context.Context, suite *TestSuite) error
	Update(ctx context.Context, id int64, name string) (*TestSuite, error)
	Delete(ctx context.Context, id int64) error
}

// TestCaseRepository defines the data access contract for test cases
type TestCaseRepository interface {
	GetByID(ctx context.Context, id int64) (*TestCase, error)
	GetAllBySuiteID(ctx context.Context, suiteID int64) ([]*TestCase, error)
	GetByName(ctx context.Context, suiteID int64, name string) (*TestCase, error)
	Create(ctx context.Context, tc *TestCase) error
	Update(ctx context.Context, id int64, name, classname string) (*TestCase, error)
	Delete(ctx context.Context, id int64) error
}

// BuildExecutionRepository defines the data access contract for build executions
type BuildExecutionRepository interface {
	GetByBuildID(ctx context.Context, buildID int64) ([]*BuildExecution, error)
	Create(ctx context.Context, execution *BuildExecution) error
	CreateBatch(ctx context.Context, executions []*BuildExecution) error
}

// FailureRepository defines the data access contract for failures
type FailureRepository interface {
	GetByID(ctx context.Context, id int64) (*Failure, error)
	GetByExecutionID(ctx context.Context, executionID int64) (*Failure, error)
	Create(ctx context.Context, failure *Failure) error
	Update(ctx context.Context, id int64, failure *Failure) (*Failure, error)
	Delete(ctx context.Context, id int64) error
}

// =============================================================================
// DOMAIN CONSTANTS AND TYPES
// =============================================================================

// Test status constants
const (
	StatusPassed  = "passed"
	StatusFailed  = "failed"
	StatusSkipped = "skipped"
	StatusError   = "error"
)

// Domain errors
type DomainError struct {
	Code    string
	Message string
}

func (e DomainError) Error() string {
	return e.Message
}

var (
	ErrProjectNotFound    = DomainError{Code: "PROJECT_NOT_FOUND", Message: "project not found"}
	ErrTestSuiteNotFound  = DomainError{Code: "TEST_SUITE_NOT_FOUND", Message: "test suite not found"}
	ErrTestCaseNotFound   = DomainError{Code: "TEST_CASE_NOT_FOUND", Message: "test case not found"}
	ErrBuildNotFound      = DomainError{Code: "BUILD_NOT_FOUND", Message: "build not found"}
	ErrInvalidInput       = DomainError{Code: "INVALID_INPUT", Message: "invalid input"}
	ErrDuplicateProject   = DomainError{Code: "DUPLICATE_PROJECT", Message: "project with this name already exists"}
	ErrDuplicateTestSuite = DomainError{Code: "DUPLICATE_TEST_SUITE", Message: "test suite with this name already exists"}
	ErrDuplicateTestCase  = errors.New("duplicate test case")
	ErrExecutionNotFound  = errors.New("execution not found")
	ErrFailureNotFound    = errors.New("failure not found")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserConfigNotFound = errors.New("user config not found")
	ErrDuplicateUser      = errors.New("user with this username already exists")
)

// =============================================================================
// FAILURE DOMAIN
// =============================================================================

type FailureService interface {
	GetFailureByID(ctx context.Context, id int64) (*Failure, error)
	GetFailureByExecutionID(ctx context.Context, executionID int64) (*Failure, error)
	CreateFailure(ctx context.Context, failure *Failure) (*Failure, error)
	UpdateFailure(ctx context.Context, id int64, failure *Failure) (*Failure, error)
	DeleteFailure(ctx context.Context, id int64) error
}

// =============================================================================
// USER DOMAIN
// =============================================================================

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type UserConfig struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	Layouts        string    `json:"layouts"` // Storing as a JSON string
	ActiveLayoutID string    `json:"active_layout_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type UserRepository interface {
	GetByID(ctx context.Context, id int) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, id int, user *User) (*User, error)
	Delete(ctx context.Context, id int) error
}

type UserConfigRepository interface {
	GetByUserID(ctx context.Context, userID int) (*UserConfig, error)
	Create(ctx context.Context, config *UserConfig) error
	Update(ctx context.Context, userID int, config *UserConfig) (*UserConfig, error)
	Delete(ctx context.Context, userID int) error
}

type UserService interface {
	GetUserByID(ctx context.Context, id int) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	CreateUser(ctx context.Context, username string) (*User, error)
	UpdateUser(ctx context.Context, id int, username string) (*User, error)
	DeleteUser(ctx context.Context, id int) error
}

type UserConfigService interface {
	GetUserConfig(ctx context.Context, userID int) (*UserConfig, error)
	CreateUserConfig(ctx context.Context, userID int, layouts, activeLayoutID string) (*UserConfig, error)
	UpdateUserConfig(ctx context.Context, userID int, layouts, activeLayoutID string) (*UserConfig, error)
	DeleteUserConfig(ctx context.Context, userID int) error
}

// =============================================================================
// SEARCH DOMAIN
// =============================================================================

type SearchResult struct {
	Type string `json:"type"`
	ID   int64  `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
}

type SearchService interface {
	Search(ctx context.Context, query string) ([]*SearchResult, error)
}

type SearchRepository interface {
	Search(ctx context.Context, query string) ([]*SearchResult, error)
}
