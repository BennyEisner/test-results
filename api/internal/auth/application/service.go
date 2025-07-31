package application

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/BennyEisner/test-results/internal/auth/domain/errors"
	"github.com/BennyEisner/test-results/internal/auth/domain/models"
	"github.com/BennyEisner/test-results/internal/auth/domain/ports"
	"github.com/markbates/goth"
)

// AuthService implements the authentication business logic
type AuthService struct {
	authRepo ports.AuthRepository
}

// NewAuthService creates a new AuthService
func NewAuthService(authRepo ports.AuthRepository) ports.AuthService {
	return &AuthService{
		authRepo: authRepo,
	}
}

// BeginOAuth2Auth starts the OAuth2 authentication flow
func (s *AuthService) BeginOAuth2Auth(ctx context.Context, provider string, state string) (string, error) {
	// Get the provider from Goth
	gothProvider, err := goth.GetProvider(provider)
	if err != nil {
		return "", errors.ErrOAuth2ProviderNotFound
	}

	// Begin authentication with Goth
	session, err := gothProvider.BeginAuth(state)
	if err != nil {
		return "", fmt.Errorf("failed to begin OAuth2 auth: %w", err)
	}

	// Get the authorization URL
	authURL, err := session.GetAuthURL()
	if err != nil {
		return "", fmt.Errorf("failed to get auth URL: %w", err)
	}

	return authURL, nil
}

// CompleteOAuth2Auth completes the OAuth2 authentication flow
func (s *AuthService) CompleteOAuth2Auth(ctx context.Context, provider string, code string, state string) (*models.User, error) {
	// Get the provider from Goth to validate it exists
	_, err := goth.GetProvider(provider)
	if err != nil {
		return nil, errors.ErrOAuth2ProviderNotFound
	}

	// For now, we'll implement a simplified version that works with our HTTP handlers
	// The actual OAuth2 completion will be handled in the HTTP layer where we have access to the request
	// This method will be called after the OAuth2 flow is completed in the HTTP handler

	// Return an error indicating this should be handled at the HTTP layer
	return nil, fmt.Errorf("OAuth2 completion should be handled at HTTP layer with full request context")
}

// CreateSession creates a new authentication session
func (s *AuthService) CreateSession(ctx context.Context, userID int64, provider string) (*models.AuthSession, error) {
	// Generate a random session ID
	sessionID, err := generateRandomString(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Create session with 24-hour expiration
	session := &models.AuthSession{
		ID:        sessionID,
		UserID:    userID,
		Provider:  provider,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		CreatedAt: time.Now(),
	}

	if err := s.authRepo.CreateSession(ctx, session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// ValidateSession validates an existing session
func (s *AuthService) ValidateSession(ctx context.Context, sessionID string) (*models.AuthContext, error) {
	session, err := s.authRepo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, errors.ErrSessionNotFound
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		// Clean up expired session
		_ = s.authRepo.DeleteSession(ctx, sessionID)
		return nil, errors.ErrSessionExpired
	}

	return &models.AuthContext{
		UserID:   session.UserID,
		Provider: session.Provider,
		IsAPIKey: false,
	}, nil
}

// DeleteSession deletes an authentication session
func (s *AuthService) DeleteSession(ctx context.Context, sessionID string) error {
	return s.authRepo.DeleteSession(ctx, sessionID)
}

// CreateAPIKey creates a new API key for a user
func (s *AuthService) CreateAPIKey(ctx context.Context, userID int64, name string) (*models.APIKey, string, error) {
	// Generate a random API key
	apiKeyBytes := make([]byte, 32)
	if _, err := rand.Read(apiKeyBytes); err != nil {
		return nil, "", fmt.Errorf("failed to generate API key: %w", err)
	}

	apiKey := base64.URLEncoding.EncodeToString(apiKeyBytes)

	// Hash the API key for storage
	keyHash := hashAPIKey(apiKey)

	// Create API key with 1-year expiration
	apiKeyModel := &models.APIKey{
		UserID:    userID,
		Name:      name,
		KeyHash:   keyHash,
		ExpiresAt: time.Now().AddDate(1, 0, 0), // 1 year
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.authRepo.CreateAPIKey(ctx, apiKeyModel); err != nil {
		return nil, "", fmt.Errorf("failed to create API key: %w", err)
	}

	return apiKeyModel, apiKey, nil
}

// ValidateAPIKey validates an API key
func (s *AuthService) ValidateAPIKey(ctx context.Context, apiKey string) (*models.AuthContext, error) {
	keyHash := hashAPIKey(apiKey)

	apiKeyModel, err := s.authRepo.GetAPIKeyByHash(ctx, keyHash)
	if err != nil {
		return nil, errors.ErrInvalidAPIKey
	}

	// Check if API key is expired
	if time.Now().After(apiKeyModel.ExpiresAt) {
		return nil, errors.ErrAPIKeyExpired
	}

	// Update last used timestamp
	_ = s.authRepo.UpdateAPIKeyLastUsed(ctx, apiKeyModel.ID)

	return &models.AuthContext{
		UserID:   apiKeyModel.UserID,
		Provider: "api_key",
		IsAPIKey: true,
		APIKeyID: apiKeyModel.ID,
	}, nil
}

// DeleteAPIKey deletes an API key
func (s *AuthService) DeleteAPIKey(ctx context.Context, userID int64, keyID int64) error {
	return s.authRepo.DeleteAPIKey(ctx, userID, keyID)
}

// ListAPIKeys lists all API keys for a user
func (s *AuthService) ListAPIKeys(ctx context.Context, userID int64) ([]*models.APIKey, error) {
	return s.authRepo.ListAPIKeysByUser(ctx, userID)
}

// GetUserByID retrieves a user by ID
func (s *AuthService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	return s.authRepo.GetUserByID(ctx, userID)
}

// GetUserByProviderID retrieves a user by provider and provider ID
func (s *AuthService) GetUserByProviderID(ctx context.Context, provider, providerID string) (*models.User, error) {
	return s.authRepo.GetUserByProviderID(ctx, provider, providerID)
}

// CreateOrUpdateUser creates or updates a user from Goth user data
func (s *AuthService) CreateOrUpdateUser(ctx context.Context, gothUser goth.User) (*models.User, error) {
	// Check if user already exists
	existingUser, err := s.authRepo.GetUserByProviderID(ctx, gothUser.Provider, gothUser.UserID)
	if err == nil && existingUser != nil {
		// Update existing user
		existingUser.Email = gothUser.Email
		existingUser.Name = gothUser.Name
		existingUser.FirstName = gothUser.FirstName
		existingUser.LastName = gothUser.LastName
		existingUser.AvatarURL = gothUser.AvatarURL
		existingUser.AccessToken = gothUser.AccessToken
		existingUser.RefreshToken = gothUser.RefreshToken
		existingUser.ExpiresAt = gothUser.ExpiresAt
		existingUser.UpdatedAt = time.Now()

		if err := s.authRepo.UpdateUser(ctx, existingUser); err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
		return existingUser, nil
	}

	// Create new user
	user := &models.User{
		Provider:     gothUser.Provider,
		ProviderID:   gothUser.UserID,
		Email:        gothUser.Email,
		Name:         gothUser.Name,
		FirstName:    gothUser.FirstName,
		LastName:     gothUser.LastName,
		AvatarURL:    gothUser.AvatarURL,
		AccessToken:  gothUser.AccessToken,
		RefreshToken: gothUser.RefreshToken,
		ExpiresAt:    gothUser.ExpiresAt,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.authRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// Helper functions

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func hashAPIKey(apiKey string) string {
	hash := sha256.Sum256([]byte(apiKey))
	return base64.URLEncoding.EncodeToString(hash[:])
}
