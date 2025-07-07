package domain

import "errors"

var (
	ErrInvalidBuildData     = errors.New("invalid build data")
	ErrBuildNotFound        = errors.New("build not found")
	ErrInvalidProjectName   = errors.New("invalid project name")
	ErrInvalidTestSuiteName = errors.New("invalid test suite name")
)
