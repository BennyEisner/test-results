package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain"
)

type UserConfigHandler struct {
	service domain.UserConfigService
}

func NewUserConfigHandler(service domain.UserConfigService) *UserConfigHandler {
	return &UserConfigHandler{service: service}
}

func (h *UserConfigHandler) GetUserConfig(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	if userIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing user ID")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	config, err := h.service.GetUserConfig(r.Context(), userID)
	if err != nil {
		switch err {
		case domain.ErrUserConfigNotFound:
			respondWithError(w, http.StatusNotFound, "user config not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, config)
}

func (h *UserConfigHandler) CreateUserConfig(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	if userIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing user ID")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	var request struct {
		Layouts        string `json:"layouts"`
		ActiveLayoutID string `json:"active_layout_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	config, err := h.service.CreateUserConfig(r.Context(), userID, request.Layouts, request.ActiveLayoutID)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to create user config")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, config)
}

func (h *UserConfigHandler) UpdateUserConfig(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	if userIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing user ID")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	var request struct {
		Layouts        string `json:"layouts"`
		ActiveLayoutID string `json:"active_layout_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	config, err := h.service.UpdateUserConfig(r.Context(), userID, request.Layouts, request.ActiveLayoutID)
	if err != nil {
		switch err {
		case domain.ErrUserConfigNotFound:
			respondWithError(w, http.StatusNotFound, "user config not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to update user config")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, config)
}

func (h *UserConfigHandler) DeleteUserConfig(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.PathValue("userID")
	if userIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing user ID")
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	if err := h.service.DeleteUserConfig(r.Context(), userID); err != nil {
		switch err {
		case domain.ErrUserConfigNotFound:
			respondWithError(w, http.StatusNotFound, "user config not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to delete user config")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "user config deleted successfully"})
}
