package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/BennyEisner/test-results/internal/auth/domain/models"
	"github.com/BennyEisner/test-results/internal/auth/domain/ports"
)

// AuthMiddleware provides authentication for protected routes
type AuthMiddleware struct {
	authService ports.AuthService
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(authService ports.AuthService) *AuthMiddleware {
	return &AuthMiddleware{authService: authService}
}

// RequireAuth middleware that requires authentication
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authContext, err := m.authenticateRequest(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Add auth context to request context
		ctx := context.WithValue(r.Context(), authContextKey, authContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth middleware that makes authentication optional
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authContext, err := m.authenticateRequest(r)
		if err == nil {
			// Add auth context to request context if authentication succeeded
			ctx := context.WithValue(r.Context(), authContextKey, authContext)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// Continue without authentication
			next.ServeHTTP(w, r)
		}
	})
}

// authenticateRequest attempts to authenticate the request using session or API key
func (m *AuthMiddleware) authenticateRequest(r *http.Request) (*models.AuthContext, error) {
	// Try API key authentication first
	if authContext, err := m.authenticateAPIKey(r); err == nil {
		return authContext, nil
	}

	// Try session authentication
	if authContext, err := m.authenticateSession(r); err == nil {
		return authContext, nil
	}

	return nil, http.ErrNoCookie // Use this as a generic auth failure
}

// authenticateAPIKey attempts to authenticate using an API key
func (m *AuthMiddleware) authenticateAPIKey(r *http.Request) (*models.AuthContext, error) {
	// Check for API key in Authorization header
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, http.ErrNoCookie
	}

	// Check if it's a Bearer token (API key)
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, http.ErrNoCookie
	}

	apiKey := strings.TrimPrefix(authHeader, "Bearer ")
	if apiKey == "" {
		return nil, http.ErrNoCookie
	}

	// Validate API key
	authContext, err := m.authService.ValidateAPIKey(r.Context(), apiKey)
	if err != nil {
		return nil, err
	}

	return authContext, nil
}

// authenticateSession attempts to authenticate using a session cookie
func (m *AuthMiddleware) authenticateSession(r *http.Request) (*models.AuthContext, error) {
	// Get session ID from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	// Validate session
	authContext, err := m.authService.ValidateSession(r.Context(), cookie.Value)
	if err != nil {
		return nil, err
	}

	return authContext, nil
}

// GetAuthContext extracts the authentication context from the request context
func GetAuthContext(r *http.Request) (*models.AuthContext, bool) {
	authCtx := r.Context().Value(authContextKey)
	if authCtx == nil {
		return nil, false
	}

	authContext, ok := authCtx.(*models.AuthContext)
	return authContext, ok
}

// RequireAPIKey middleware that requires API key authentication (for CLI/Jenkins)
func (m *AuthMiddleware) RequireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authContext, err := m.authenticateAPIKey(r)
		if err != nil {
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}

		// Verify it's an API key authentication
		if !authContext.IsAPIKey {
			http.Error(w, "API key required", http.StatusUnauthorized)
			return
		}

		// Add auth context to request context
		ctx := context.WithValue(r.Context(), authContextKey, authContext)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// contextKey is a private type for context keys in this package
type contextKey string

const authContextKey contextKey = "auth_context"
