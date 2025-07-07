package domain

import "errors"

var (
	ErrInvalidExecutionData   = errors.New("invalid execution data")
	ErrBuildExecutionNotFound = errors.New("build execution not found")
	ErrInvalidBuildData       = errors.New("invalid build data")
)
