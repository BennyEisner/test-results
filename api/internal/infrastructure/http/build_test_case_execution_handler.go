package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// BuildTestCaseExecutionHandler handles HTTP requests for build test case executions
type BuildTestCaseExecutionHandler struct {
	Service ports.BuildTestCaseExecutionService
}

// NewBuildTestCaseExecutionHandler creates a new BuildTestCaseExecutionHandler
func NewBuildTestCaseExecutionHandler(service ports.BuildTestCaseExecutionService) *BuildTestCaseExecutionHandler {
	return &BuildTestCaseExecutionHandler{Service: service}
}

// GetExecutionByID handles GET /executions/{id}
func (h *BuildTestCaseExecutionHandler) GetExecutionByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	execution, err := h.Service.GetExecutionByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if execution == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(execution); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetExecutionsByBuildID handles GET /builds/{buildID}/executions
func (h *BuildTestCaseExecutionHandler) GetExecutionsByBuildID(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.URL.Query().Get("build_id")
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid build_id", http.StatusBadRequest)
		return
	}
	executions, err := h.Service.GetExecutionsByBuildID(r.Context(), buildID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(executions); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateExecution handles POST /builds/{buildID}/executions
func (h *BuildTestCaseExecutionHandler) CreateExecution(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.URL.Query().Get("build_id")
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid build_id", http.StatusBadRequest)
		return
	}
	var input models.BuildExecutionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	execution, err := h.Service.CreateExecution(r.Context(), buildID, &input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(execution); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateExecution handles PUT /executions/{id}
func (h *BuildTestCaseExecutionHandler) UpdateExecution(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var execution models.BuildTestCaseExecution
	if err := json.NewDecoder(r.Body).Decode(&execution); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	updatedExecution, err := h.Service.UpdateExecution(r.Context(), id, &execution)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if updatedExecution == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(updatedExecution); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteExecution handles DELETE /executions/{id}
func (h *BuildTestCaseExecutionHandler) DeleteExecution(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteExecution(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
