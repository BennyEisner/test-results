package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain"
)

type UserHandler struct {
	service domain.UserService
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing user ID")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	user, err := h.service.GetUserByID(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, "user not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		respondWithError(w, http.StatusBadRequest, "missing username parameter")
		return
	}

	user, err := h.service.GetUserByUsername(r.Context(), username)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, "user not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.CreateUser(r.Context(), request.Username)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "username is required")
		case domain.ErrDuplicateUser:
			respondWithError(w, http.StatusConflict, "user with this username already exists")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to create user")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing user ID")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	var request struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.service.UpdateUser(r.Context(), id, request.Username)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, "user not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "username is required")
		case domain.ErrDuplicateUser:
			respondWithError(w, http.StatusConflict, "user with this username already exists")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to update user")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing user ID")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID format")
		return
	}

	if err := h.service.DeleteUser(r.Context(), id); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			respondWithError(w, http.StatusNotFound, "user not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to delete user")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "user deleted successfully"})
}
