package http

import (
	"encoding/json"
	"net/http"

	"github.com/BennyEisner/test-results/internal/auth/infrastructure/middleware"
	"github.com/BennyEisner/test-results/internal/user_config/domain/models"
	"github.com/BennyEisner/test-results/internal/user_config/domain/ports"
)

// UserConfigHandler handles HTTP requests for user configurations
type UserConfigHandler struct {
	Service ports.UserConfigService
}

// NewUserConfigHandler creates a new UserConfigHandler
func NewUserConfigHandler(service ports.UserConfigService) *UserConfigHandler {
	return &UserConfigHandler{Service: service}
}

// GetUserConfigs handles GET /configs
// @Summary Get user configs
// @Description Retrieve all configuration settings for the authenticated user
// @Tags user-configs
// @Accept json
// @Produce json
// @Success 200 {array} object
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /configs [get]
func (h *UserConfigHandler) GetUserConfigs(w http.ResponseWriter, r *http.Request) {
	authContext, ok := middleware.GetAuthContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	configs, err := h.Service.GetUserConfigs(r.Context(), authContext.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(configs); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// SaveUserConfig handles POST /configs
// @Summary Save user config
// @Description Create or update a configuration setting for the authenticated user
// @Tags user-configs
// @Accept json
// @Produce json
// @Param config body object true "Config creation request"
// @Success 201 {object} object
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /configs [post]
func (h *UserConfigHandler) SaveUserConfig(w http.ResponseWriter, r *http.Request) {
	authContext, ok := middleware.GetAuthContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.UserConfig
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	config, err := h.Service.SaveUserConfig(
		r.Context(),
		authContext.UserID,
		req.Layouts,
		req.ActiveLayoutID,
	)
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

// UpdateActiveLayoutID handles PUT /configs/active
// @Summary Update active layout ID
// @Description Update only the active layout ID for the authenticated user's configuration
// @Tags user-configs
// @Accept json
// @Produce json
// @Param request body object true "Active layout ID update request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /configs/active [put]
func (h *UserConfigHandler) UpdateActiveLayoutID(w http.ResponseWriter, r *http.Request) {
	authContext, ok := middleware.GetAuthContext(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		ActiveLayoutID string `json:"active_layout_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := h.Service.UpdateActiveLayoutID(
		r.Context(),
		authContext.UserID,
		req.ActiveLayoutID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{"message": "Active layout ID updated successfully"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
