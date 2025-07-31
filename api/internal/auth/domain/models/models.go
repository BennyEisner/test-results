package models

import (
	"time"
)

// User represents an authenticated user in the system
type User struct {
	ID           int64     `json:"id"`
	Provider     string    `json:"provider"`    // "okta", "github", etc.
	ProviderID   string    `json:"provider_id"` // External provider's user ID
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	AvatarURL    string    `json:"avatar_url"`
	AccessToken  string    `json:"-"` // Not exposed in JSON
	RefreshToken string    `json:"-"` // Not exposed in JSON
	ExpiresAt    time.Time `json:"expires_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// APIKey represents an API key for CLI/Jenkins authentication
type APIKey struct {
	ID         int64      `json:"id"`
	UserID     int64      `json:"user_id"`
	Name       string     `json:"name"` // Human-readable name for the key
	KeyHash    string     `json:"-"`    // Hashed API key
	LastUsedAt *time.Time `json:"last_used_at"`
	ExpiresAt  time.Time  `json:"expires_at"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// AuthSession represents an active authentication session
type AuthSession struct {
	ID        string    `json:"id"`
	UserID    int64     `json:"user_id"`
	Provider  string    `json:"provider"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthContext represents the current authentication context
type AuthContext struct {
	UserID   int64  `json:"user_id"`
	Provider string `json:"provider"`
	IsAPIKey bool   `json:"is_api_key"`
	APIKeyID int64  `json:"api_key_id,omitempty"`
}
