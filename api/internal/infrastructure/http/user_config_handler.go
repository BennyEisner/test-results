package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// UserConfigHandler handles HTTP requests for user configurations
type UserConfigHandler struct {
	Service ports.UserConfigService
}

// NewUserConfigHandler creates a new UserConfigHandler
func NewUserConfigHandler(service ports.UserConfigService) *UserConfigHandler {
	return &UserConfigHandler{Service: service}
}

// GetUserConfig handles GET /users/{userID}/config
func (h *UserConfigHandler) GetUserConfig(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	config, err := h.Service.GetUserConfig(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if config == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateUserConfig handles POST /users/{userID}/config
func (h *UserConfigHandler) CreateUserConfig(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	var req struct {
		Layouts        string `json:"layouts"`
		ActiveLayoutID string `json:"active_layout_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	config, err := h.Service.CreateUserConfig(r.Context(), userID, req.Layouts, req.ActiveLayoutID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateUserConfig handles PUT /users/{userID}/config
func (h *UserConfigHandler) UpdateUserConfig(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	var req struct {
		Layouts        string `json:"layouts"`
		ActiveLayoutID string `json:"active_layout_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	config, err := h.Service.UpdateUserConfig(r.Context(), userID, req.Layouts, req.ActiveLayoutID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if config == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(config); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteUserConfig handles DELETE /users/{userID}/config
func (h *UserConfigHandler) DeleteUserConfig(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteUserConfig(r.Context(), userID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
