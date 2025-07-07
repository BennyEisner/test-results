package errors

import "errors"

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserExists     = errors.New("user already exists")
	ErrInvalidUser    = errors.New("invalid user data")
	ErrConfigNotFound = errors.New("user config not found")
	ErrInvalidConfig  = errors.New("invalid user config data")
)
