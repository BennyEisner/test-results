package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// BuildHandler handles HTTP requests for builds
type BuildHandler struct {
	Service ports.BuildService
}

// NewBuildHandler creates a new BuildHandler
func NewBuildHandler(service ports.BuildService) *BuildHandler {
	return &BuildHandler{Service: service}
}

// GetBuildByID handles GET /builds/{id}
func (h *BuildHandler) GetBuildByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	build, err := h.Service.GetBuildByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if build == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(build); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetBuildsByProjectID handles GET /projects/{projectID}/builds
func (h *BuildHandler) GetBuildsByProjectID(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
		return
	}
	builds, err := h.Service.GetBuildsByProjectID(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(builds); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetBuildsByTestSuiteID handles GET /test-suites/{suiteID}/builds
func (h *BuildHandler) GetBuildsByTestSuiteID(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.URL.Query().Get("suite_id")
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid suite_id", http.StatusBadRequest)
		return
	}
	builds, err := h.Service.GetBuildsByTestSuiteID(r.Context(), suiteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(builds); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateBuild handles POST /builds
func (h *BuildHandler) CreateBuild(w http.ResponseWriter, r *http.Request) {
	var build models.Build
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	createdBuild, err := h.Service.CreateBuild(r.Context(), &build)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(createdBuild); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateBuild handles PUT /builds/{id}
func (h *BuildHandler) UpdateBuild(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var build models.Build
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	updatedBuild, err := h.Service.UpdateBuild(r.Context(), id, &build)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if updatedBuild == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(updatedBuild); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteBuild handles DELETE /builds/{id}
func (h *BuildHandler) DeleteBuild(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteBuild(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
