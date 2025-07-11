package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BennyEisner/test-results/internal/auth/domain/models"
	"github.com/BennyEisner/test-results/internal/auth/domain/ports"
	"github.com/markbates/goth/gothic"
)

// AuthHandler handles HTTP requests for authentication
type AuthHandler struct {
	authService ports.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService ports.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// BeginOAuth2Auth handles the start of OAuth2 authentication
// @Summary Begin OAuth2 authentication
// @Description Start OAuth2 authentication flow with the specified provider
// @Tags auth
// @Accept json
// @Produce json
// @Param provider path string true "OAuth2 provider name"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/{provider} [get]
func (h *AuthHandler) BeginOAuth2Auth(w http.ResponseWriter, r *http.Request) {
	// Extract provider from URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "invalid provider path", http.StatusBadRequest)
		return
	}
	provider := pathParts[2]

	// Generate state for CSRF protection
	state := r.URL.Query().Get("state")
	if state == "" {
		// Generate random state if not provided
		state = generateRandomState()
	}

	// Begin OAuth2 authentication
	authURL, err := h.authService.BeginOAuth2Auth(r.Context(), provider, state)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Redirect to OAuth2 provider
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// OAuth2Callback handles the OAuth2 callback
// @Summary OAuth2 callback
// @Description Handle OAuth2 callback from provider
// @Tags auth
// @Accept json
// @Produce json
// @Param provider path string true "OAuth2 provider name"
// @Success 200 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/{provider}/callback [get]
func (h *AuthHandler) OAuth2Callback(w http.ResponseWriter, r *http.Request) {
	// Extract provider from URL path (for validation)
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "invalid callback path", http.StatusBadRequest)
		return
	}
	_ = pathParts[2] // provider - used for validation but not needed for Goth

	// Complete OAuth2 authentication using Goth
	gothUser, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		http.Error(w, "authentication failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create or update user in our database
	user, err := h.authService.CreateOrUpdateUser(r.Context(), gothUser)
	if err != nil {
		http.Error(w, "failed to create/update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Create session
	session, err := h.authService.CreateSession(r.Context(), user.ID, user.Provider)
	if err != nil {
		http.Error(w, "failed to create session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Set session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    session.ID,
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil, // Secure in production
		MaxAge:   86400,        // 24 hours
	})

	// Redirect to frontend with success
	http.Redirect(w, r, "/auth/success", http.StatusTemporaryRedirect)
}

// Logout handles user logout
// @Summary Logout user
// @Description Logout user and clear session
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/logout [post]
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get session ID from cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No session cookie, already logged out
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
		return
	}

	// Delete session
	err = h.authService.DeleteSession(r.Context(), cookie.Value)
	if err != nil {
		// Log error but don't fail the request
		// The session might already be expired
	}

	// Clear session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		MaxAge:   -1, // Delete cookie
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
}

// GetCurrentUser returns the current authenticated user
// @Summary Get current user
// @Description Get information about the currently authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get auth context from middleware
	authCtx := r.Context().Value("auth_context")
	if authCtx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authContext := authCtx.(*models.AuthContext)

	// Get user details
	user, err := h.authService.GetUserByID(r.Context(), authContext.UserID)
	if err != nil {
		http.Error(w, "failed to get user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// CreateAPIKey creates a new API key for the authenticated user
// @Summary Create API key
// @Description Create a new API key for CLI/Jenkins authentication
// @Tags auth
// @Accept json
// @Produce json
// @Param request body CreateAPIKeyRequest true "API key creation request"
// @Success 200 {object} CreateAPIKeyResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/api-keys [post]
func (h *AuthHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	// Get auth context from middleware
	authCtx := r.Context().Value("auth_context")
	if authCtx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authContext := authCtx.(*models.AuthContext)

	// Parse request
	var req CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	// Create API key
	apiKey, plainTextKey, err := h.authService.CreateAPIKey(r.Context(), authContext.UserID, req.Name)
	if err != nil {
		http.Error(w, "failed to create API key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := CreateAPIKeyResponse{
		APIKey:       apiKey,
		PlainTextKey: plainTextKey,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListAPIKeys lists all API keys for the authenticated user
// @Summary List API keys
// @Description List all API keys for the authenticated user
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {array} models.APIKey
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/api-keys [get]
func (h *AuthHandler) ListAPIKeys(w http.ResponseWriter, r *http.Request) {
	// Get auth context from middleware
	authCtx := r.Context().Value("auth_context")
	if authCtx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authContext := authCtx.(*models.AuthContext)

	// List API keys
	apiKeys, err := h.authService.ListAPIKeys(r.Context(), authContext.UserID)
	if err != nil {
		http.Error(w, "failed to list API keys: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiKeys)
}

// DeleteAPIKey deletes an API key
// @Summary Delete API key
// @Description Delete an API key by ID
// @Tags auth
// @Accept json
// @Produce json
// @Param id path int true "API key ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/api-keys/{id} [delete]
func (h *AuthHandler) DeleteAPIKey(w http.ResponseWriter, r *http.Request) {
	// Get auth context from middleware
	authCtx := r.Context().Value("auth_context")
	if authCtx == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authContext := authCtx.(*models.AuthContext)

	// Extract API key ID from URL
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 4 {
		http.Error(w, "invalid API key ID", http.StatusBadRequest)
		return
	}

	keyID, err := strconv.ParseInt(pathParts[3], 10, 64)
	if err != nil {
		http.Error(w, "invalid API key ID", http.StatusBadRequest)
		return
	}

	// Delete API key
	err = h.authService.DeleteAPIKey(r.Context(), authContext.UserID, keyID)
	if err != nil {
		http.Error(w, "failed to delete API key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "API key deleted"})
}

// Request/Response types

type CreateAPIKeyRequest struct {
	Name string `json:"name"`
}

type CreateAPIKeyResponse struct {
	APIKey       *models.APIKey `json:"api_key"`
	PlainTextKey string         `json:"plain_text_key"`
}

// Helper functions

func generateRandomState() string {
	// This is a simplified version - in production, use crypto/rand
	return "state_" + strconv.FormatInt(time.Now().UnixNano(), 10)
}
