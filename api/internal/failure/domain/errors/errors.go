package errors

import "errors"

var (
	ErrFailureNotFound = errors.New("failure not found")
	ErrInvalidFailure  = errors.New("invalid failure data")
)
