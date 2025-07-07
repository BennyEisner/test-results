package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/user/domain/ports"
)

// UserConfigHandler handles HTTP requests for user configs
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
		respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	ctx := r.Context()
	configs, err := h.Service.GetUserConfigs(ctx, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, configs)
}

// GetUserConfig handles GET /users/{userID}/configs/{key}
// @Summary Get user config by key
// @Description Retrieve a specific configuration setting for a user by key
// @Tags user-configs
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param key path string true "Config key"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{userID}/configs/{key} [get]
func (h *UserConfigHandler) GetUserConfig(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	key := r.PathValue("key")
	if key == "" {
		respondWithError(w, http.StatusBadRequest, "missing config key")
		return
	}

	ctx := r.Context()
	config, err := h.Service.GetUserConfig(ctx, userID, key)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, config)
}

// SetUserConfig handles POST /users/{userID}/configs
// @Summary Set user config
// @Description Create or update a configuration setting for a user
// @Tags user-configs
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param config body object true "Config creation request" schema="{key:string,value:string}"
// @Success 201 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{userID}/configs [post]
func (h *UserConfigHandler) SetUserConfig(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	config, err := h.Service.SetUserConfig(ctx, userID, req.Key, req.Value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, config)
}

// UpdateUserConfig handles PUT /configs/{id}
// @Summary Update user config
// @Description Update an existing configuration setting by ID
// @Tags user-configs
// @Accept json
// @Produce json
// @Param id path int true "Config ID"
// @Param config body object true "Config update request" schema="{value:string}"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /configs/{id} [put]
func (h *UserConfigHandler) UpdateUserConfig(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid config ID")
		return
	}

	var req struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	config, err := h.Service.UpdateUserConfig(ctx, id, req.Value)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, config)
}

// DeleteUserConfig handles DELETE /configs/{id}
// @Summary Delete user config
// @Description Delete a configuration setting by its ID
// @Tags user-configs
// @Accept json
// @Produce json
// @Param id path int true "Config ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /configs/{id} [delete]
func (h *UserConfigHandler) DeleteUserConfig(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid config ID")
		return
	}

	ctx := r.Context()
	if err := h.Service.DeleteUserConfig(ctx, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
