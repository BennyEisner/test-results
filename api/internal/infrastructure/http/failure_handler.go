package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// FailureHandler handles HTTP requests for failures
type FailureHandler struct {
	Service ports.FailureService
}

// NewFailureHandler creates a new FailureHandler
func NewFailureHandler(service ports.FailureService) *FailureHandler {
	return &FailureHandler{Service: service}
}

// GetFailureByID handles GET /failures/{id}
func (h *FailureHandler) GetFailureByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	failure, err := h.Service.GetFailureByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if failure == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(failure); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetFailureByExecutionID handles GET /executions/{executionID}/failure
func (h *FailureHandler) GetFailureByExecutionID(w http.ResponseWriter, r *http.Request) {
	executionIDStr := r.URL.Query().Get("execution_id")
	executionID, err := strconv.ParseInt(executionIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid execution_id", http.StatusBadRequest)
		return
	}
	failure, err := h.Service.GetFailureByExecutionID(r.Context(), executionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if failure == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(failure); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateFailure handles POST /failures
func (h *FailureHandler) CreateFailure(w http.ResponseWriter, r *http.Request) {
	var failure models.Failure
	if err := json.NewDecoder(r.Body).Decode(&failure); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	createdFailure, err := h.Service.CreateFailure(r.Context(), &failure)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdFailure); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateFailure handles PUT /failures/{id}
func (h *FailureHandler) UpdateFailure(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var failure models.Failure
	if err := json.NewDecoder(r.Body).Decode(&failure); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	updatedFailure, err := h.Service.UpdateFailure(r.Context(), id, &failure)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if updatedFailure == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(updatedFailure); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteFailure handles DELETE /failures/{id}
func (h *FailureHandler) DeleteFailure(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteFailure(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
