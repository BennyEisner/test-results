package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/auth/domain/models"
	"github.com/markbates/goth"
)

// AuthService defines the authentication business logic
type AuthService interface {
	// OAuth2 authentication
	BeginOAuth2Auth(ctx context.Context, provider string, state string) (string, error)
	CompleteOAuth2Auth(ctx context.Context, provider string, code string, state string) (*models.User, error)

	// Session management
	CreateSession(ctx context.Context, userID int64, provider string) (*models.AuthSession, error)
	ValidateSession(ctx context.Context, sessionID string) (*models.AuthContext, error)
	DeleteSession(ctx context.Context, sessionID string) error

	// API Key management
	CreateAPIKey(ctx context.Context, userID int64, name string) (*models.APIKey, string, error) // Returns key and plain text
	ValidateAPIKey(ctx context.Context, apiKey string) (*models.AuthContext, error)
	DeleteAPIKey(ctx context.Context, userID int64, keyID int64) error
	ListAPIKeys(ctx context.Context, userID int64) ([]*models.APIKey, error)

	// User management
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	GetUserByProviderID(ctx context.Context, provider, providerID string) (*models.User, error)
	CreateOrUpdateUser(ctx context.Context, gothUser goth.User) (*models.User, error)
}

// AuthRepository defines the data access interface for authentication
type AuthRepository interface {
	// User operations
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	GetUserByProviderID(ctx context.Context, provider, providerID string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	UpsertUser(ctx context.Context, user *models.User) error

	// Session operations
	CreateSession(ctx context.Context, session *models.AuthSession) error
	GetSession(ctx context.Context, sessionID string) (*models.AuthSession, error)
	DeleteSession(ctx context.Context, sessionID string) error
	DeleteExpiredSessions(ctx context.Context) error

	// API Key operations
	CreateAPIKey(ctx context.Context, apiKey *models.APIKey) error
	GetAPIKeyByHash(ctx context.Context, keyHash string) (*models.APIKey, error)
	UpdateAPIKeyLastUsed(ctx context.Context, keyID int64) error
	DeleteAPIKey(ctx context.Context, keyID int64) error
	ListAPIKeysByUser(ctx context.Context, userID int64) ([]*models.APIKey, error)
}
