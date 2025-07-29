package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/build/application"
	"github.com/BennyEisner/test-results/internal/build/domain/models"
)

// BuildHandler handles HTTP requests for builds
type BuildHandler struct {
	Service application.BuildService
}

// NewBuildHandler creates a new BuildHandler
func NewBuildHandler(service application.BuildService) *BuildHandler {
	return &BuildHandler{Service: service}
}

// GetBuildByID handles GET /builds/{id}
// @Summary Get build by ID
// @Description Retrieve a build by its unique identifier
// @Tags builds
// @Accept json
// @Produce json
// @Param id path int true "Build ID"
// @Success 200 {object} models.Build
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /builds/{id} [get]
func (h *BuildHandler) GetBuildByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid build ID")
		return
	}
	ctx := r.Context()
	build, err := h.Service.GetBuildByID(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, build)
}

// GetBuilds handles GET /builds
// @Summary Get builds by project or test suite
// @Description Retrieve builds by either project_id or suite_id
// @Tags builds
// @Accept json
// @Produce json
// @Param project_id query int true "Project ID"
// @Param suite_id query int false "Test Suite ID"
// @Success 200 {array} models.Build
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /builds [get]
func (h *BuildHandler) GetBuilds(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	if projectIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing project_id")
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project_id")
		return
	}

	var suiteID *int64
	suiteIDStr := r.URL.Query().Get("suite_id")
	if suiteIDStr != "" {
		id, err := strconv.ParseInt(suiteIDStr, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid suite_id")
			return
		}
		suiteID = &id
	}

	ctx := r.Context()
	builds, err := h.Service.GetBuilds(ctx, projectID, suiteID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, builds)
}

// CreateBuild handles POST /builds
// @Summary Create a new build
// @Description Create a new build for a project and test suite
// @Tags builds
// @Accept json
// @Produce json
// @Param build body models.Build true "Build to create"
// @Success 201 {object} models.Build
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /builds [post]
func (h *BuildHandler) CreateBuild(w http.ResponseWriter, r *http.Request) {
	var build models.Build
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx := r.Context()
	id, err := h.Service.CreateBuild(ctx, &build)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	build.ID = id
	respondWithJSON(w, http.StatusCreated, build)
}

// UpdateBuild handles PUT /builds/{id}
// @Summary Update a build
// @Description Update an existing build's details
// @Tags builds
// @Accept json
// @Produce json
// @Param id path int true "Build ID"
// @Param build body models.Build true "Build to update"
// @Success 200 {object} models.Build
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /builds/{id} [put]
func (h *BuildHandler) UpdateBuild(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid build ID")
		return
	}
	var build models.Build
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	build.ID = id
	ctx := r.Context()
	err = h.Service.UpdateBuild(ctx, &build)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, build)
}

// DeleteBuild handles DELETE /builds/{id}
// @Summary Delete a build
// @Description Delete a build by its ID
// @Tags builds
// @Accept json
// @Produce json
// @Param id path int true "Build ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /builds/{id} [delete]
func (h *BuildHandler) DeleteBuild(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid build ID")
		return
	}
	if err := h.Service.DeleteBuild(r.Context(), id); err != nil {
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
