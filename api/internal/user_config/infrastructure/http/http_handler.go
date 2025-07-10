package http

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// GetUserConfigs handles GET /users/{userID}/configs
// @Summary Get user configs
// @Description Retrieve all configuration settings for a specific user
// @Tags user-configs
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Success 200 {array} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{userID}/configs [get]
func (h *UserConfigHandler) GetUserConfigs(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	configs, err := h.Service.GetUserConfigs(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(configs); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// SaveUserConfig handles POST /users/{userID}/configs
// @Summary Save user config
// @Description Create or update a configuration setting for a user
// @Tags user-configs
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param config body object true "Config creation request"
// @Success 201 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{userID}/configs [post]
func (h *UserConfigHandler) SaveUserConfig(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}
	var req models.UserConfig
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	config, err := h.Service.SaveUserConfig(r.Context(), userID, req.Layouts, req.ActiveLayoutID)
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

// UpdateActiveLayoutID handles PUT /users/{userID}/configs/active
// @Summary Update active layout ID
// @Description Update only the active layout ID for a user's configuration
// @Tags user-configs
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param request body object true "Active layout ID update request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{userID}/configs/active [put]
func (h *UserConfigHandler) UpdateActiveLayoutID(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	var req struct {
		ActiveLayoutID string `json:"active_layout_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = h.Service.UpdateActiveLayoutID(r.Context(), userID, req.ActiveLayoutID)
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
