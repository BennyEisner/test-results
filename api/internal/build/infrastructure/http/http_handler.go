package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/build/domain/ports"
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
// @Summary Get build by ID
// @Description Retrieve a build by its unique identifier
// @Tags builds
// @Accept json
// @Produce json
// @Param id path int true "Build ID"
// @Success 200 {object} object
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
	build, err := h.Service.GetBuild(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, build)
}

// GetBuilds handles GET /builds
// @Summary Get builds
// @Description Retrieve builds, optionally filtered by project_id or suite_id
// @Tags builds
// @Accept json
// @Produce json
// @Param project_id query int false "Project ID"
// @Param suite_id query int false "Test Suite ID"
// @Success 200 {array} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /builds [get]
func (h *BuildHandler) GetBuilds(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	suiteIDStr := r.URL.Query().Get("suite_id")
	ctx := r.Context()

	if projectIDStr != "" {
		projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid project_id")
			return
		}
		builds, err := h.Service.GetBuildsByProject(ctx, projectID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, builds)
		return
	}

	if suiteIDStr != "" {
		suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid suite_id")
			return
		}
		builds, err := h.Service.GetBuildsByTestSuite(ctx, suiteID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, builds)
		return
	}

	respondWithError(w, http.StatusBadRequest, "project_id or suite_id must be provided")
}

// GetBuildsByTestSuite handles GET /test-suites/{suiteID}/builds
// @Summary Get builds by test suite ID
// @Description Retrieve all builds for a specific test suite
// @Tags builds
// @Accept json
// @Produce json
// @Param suite_id query int true "Test Suite ID"
// @Success 200 {array} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /builds [get]
func (h *BuildHandler) GetBuildsByTestSuite(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.URL.Query().Get("suite_id")
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid suite_id")
		return
	}
	ctx := r.Context()
	builds, err := h.Service.GetBuildsByTestSuite(ctx, suiteID)
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
// @Param build body object true "Build creation request" schema="{project_id:int,suite_id:int,build_number:string}"
// @Success 201 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /builds [post]
func (h *BuildHandler) CreateBuild(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectID   int64  `json:"project_id"`
		SuiteID     int64  `json:"suite_id"`
		BuildNumber string `json:"build_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx := r.Context()
	build, err := h.Service.CreateBuild(ctx, req.ProjectID, req.SuiteID, req.BuildNumber)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, build)
}

// UpdateBuild handles PUT /builds/{id}
// @Summary Update a build
// @Description Update an existing build's build number
// @Tags builds
// @Accept json
// @Produce json
// @Param id path int true "Build ID"
// @Param build body object true "Build update request" schema="{build_number:string}"
// @Success 200 {object} object
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
	var req struct {
		BuildNumber string `json:"build_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx := r.Context()
	build, err := h.Service.UpdateBuild(ctx, id, req.BuildNumber)
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
