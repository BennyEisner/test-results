package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/service"
)

type UserConfigHandler struct {
	service *service.UserConfigService
}

func NewUserConfigHandler(s *service.UserConfigService) *UserConfigHandler {
	return &UserConfigHandler{service: s}
}

func (h *UserConfigHandler) GetUserConfig(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	config, err := h.service.GetUserConfig(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if config == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(config)
}

type SaveConfigRequest struct {
	Layouts        interface{} `json:"layouts"`
	ActiveLayoutID string      `json:"activeLayoutId"`
}

func (h *UserConfigHandler) SaveUserConfig(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.PathValue("userId"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req SaveConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.service.SaveUserConfig(userID, req.Layouts, req.ActiveLayoutID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
