package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/user/domain/ports"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	Service ports.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(service ports.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// GetUserByID handles GET /users/{id}
// @Summary Get user by ID
// @Description Retrieve a user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	ctx := r.Context()
	user, err := h.Service.GetUser(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// GetUserByUsername handles GET /users/username/{username}
// @Summary Get user by username
// @Description Retrieve a user by their username
// @Tags users
// @Accept json
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/username/{username} [get]
func (h *UserHandler) GetUserByUsername(w http.ResponseWriter, r *http.Request) {
	username := r.PathValue("username")
	if username == "" {
		respondWithError(w, http.StatusBadRequest, "missing username")
		return
	}

	ctx := r.Context()
	user, err := h.Service.GetUserByUsername(ctx, username)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// GetUserByEmail handles GET /users/email/{email}
// @Summary Get user by email
// @Description Retrieve a user by their email address
// @Tags users
// @Accept json
// @Produce json
// @Param email path string true "Email address"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/email/{email} [get]
func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email := r.PathValue("email")
	if email == "" {
		respondWithError(w, http.StatusBadRequest, "missing email")
		return
	}

	ctx := r.Context()
	user, err := h.Service.GetUserByEmail(ctx, email)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// CreateUser handles POST /users
// @Summary Create a new user
// @Description Create a new user with username and email
// @Tags users
// @Accept json
// @Produce json
// @Param user body object true "User creation request" schema="{username:string,email:string}"
// @Success 201 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	user, err := h.Service.CreateUser(ctx, req.Username, req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// UpdateUser handles PUT /users/{id}
// @Summary Update a user
// @Description Update an existing user's username and email
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body object true "User update request" schema="{username:string,email:string}"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	user, err := h.Service.UpdateUser(ctx, id, req.Username, req.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// DeleteUser handles DELETE /users/{id}
// @Summary Delete a user
// @Description Delete a user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid user ID")
		return
	}

	ctx := r.Context()
	if err := h.Service.DeleteUser(ctx, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions for HTTP responses
func respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
