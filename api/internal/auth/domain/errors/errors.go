package errors

import "errors"

var (
	// Authentication errors
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrSessionExpired     = errors.New("session expired")
	ErrSessionNotFound    = errors.New("session not found")
	ErrInvalidAPIKey      = errors.New("invalid API key")
	ErrAPIKeyExpired      = errors.New("API key expired")
	ErrAPIKeyNotFound     = errors.New("API key not found")

	// OAuth2 errors
	ErrOAuth2ProviderNotFound = errors.New("OAuth2 provider not found")
	ErrOAuth2StateMismatch    = errors.New("OAuth2 state mismatch")
	ErrOAuth2CodeInvalid      = errors.New("OAuth2 authorization code invalid")
	ErrOAuth2TokenExchange    = errors.New("OAuth2 token exchange failed")

	// User errors
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserData   = errors.New("invalid user data")

	// Authorization errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)
