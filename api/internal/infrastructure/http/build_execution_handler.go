package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// BuildExecutionHandler handles HTTP requests for build executions
type BuildExecutionHandler struct {
	Service ports.BuildExecutionService
}

// NewBuildExecutionHandler creates a new BuildExecutionHandler
func NewBuildExecutionHandler(service ports.BuildExecutionService) *BuildExecutionHandler {
	return &BuildExecutionHandler{Service: service}
}

// GetBuildExecutions handles GET /builds/{buildID}/executions
func (h *BuildExecutionHandler) GetBuildExecutions(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.URL.Query().Get("build_id")
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid build_id", http.StatusBadRequest)
		return
	}
	executions, err := h.Service.GetBuildExecutions(r.Context(), buildID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(executions); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateBuildExecutions handles POST /builds/{buildID}/executions
func (h *BuildExecutionHandler) CreateBuildExecutions(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.URL.Query().Get("build_id")
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid build_id", http.StatusBadRequest)
		return
	}
	var executions []*models.BuildExecution
	if err := json.NewDecoder(r.Body).Decode(&executions); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if err := h.Service.CreateBuildExecutions(r.Context(), buildID, executions); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
