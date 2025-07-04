package domain

import "errors"

// Domain error constants
var (
	ErrProjectNotFound      = errors.New("project not found")
	ErrProjectAlreadyExists = errors.New("project already exists")
	ErrInvalidProjectName   = errors.New("invalid project name")

	ErrTestSuiteNotFound      = errors.New("test suite not found")
	ErrTestSuiteAlreadyExists = errors.New("test suite already exists")
	ErrDuplicateTestSuite     = errors.New("test suite already exists")
	ErrInvalidTestSuiteName   = errors.New("invalid test suite name")

	ErrTestCaseNotFound      = errors.New("test case not found")
	ErrTestCaseAlreadyExists = errors.New("test case already exists")
	ErrInvalidTestCaseName   = errors.New("invalid test case name")

	ErrBuildNotFound      = errors.New("build not found")
	ErrBuildAlreadyExists = errors.New("build already exists")
	ErrInvalidBuildData   = errors.New("invalid build data")

	ErrBuildExecutionNotFound = errors.New("build execution not found")
	ErrInvalidExecutionData   = errors.New("invalid execution data")

	ErrFailureNotFound = errors.New("failure not found")

	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUsername   = errors.New("invalid username")

	ErrUserConfigNotFound = errors.New("user config not found")
)
