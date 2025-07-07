package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/project/domain"
	"github.com/BennyEisner/test-results/internal/project/domain/ports"
)

// ProjectHandler handles HTTP requests for project operations
type ProjectHandler struct {
	projectService ports.ProjectService
}

// NewProjectHandler creates a new project handler
func NewProjectHandler(projectService ports.ProjectService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
	}
}

// GetProjectByID handles GET /api/projects/{id}
// @Summary Get project by ID
// @Description Retrieve a project by its unique identifier
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id} [get]
func (h *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing project ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID format")
		return
	}

	ctx := r.Context()
	project, err := h.projectService.GetProject(ctx, id)
	if err != nil {
		switch err {
		case domain.ErrProjectNotFound:
			respondWithError(w, http.StatusNotFound, "project not found")
		case domain.ErrInvalidProjectName:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, project)
}

// GetAllProjects handles GET /api/projects
// @Summary Get all projects
// @Description Retrieve all projects in the system
// @Tags projects
// @Accept json
// @Produce json
// @Success 200 {array} object
// @Failure 500 {object} map[string]string
// @Router /projects [get]
func (h *ProjectHandler) GetAllProjects(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projects, err := h.projectService.GetAllProjects(ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to retrieve projects")
		return
	}

	respondWithJSON(w, http.StatusOK, projects)
}

// CreateProject handles POST /api/projects
// @Summary Create a new project
// @Description Create a new project with the given name
// @Tags projects
// @Accept json
// @Produce json
// @Param project body object true "Project creation request" schema="{name:string}"
// @Success 201 {object} object
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects [post]
func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	project, err := h.projectService.CreateProject(ctx, request.Name)
	if err != nil {
		switch err {
		case domain.ErrInvalidProjectName:
			respondWithError(w, http.StatusBadRequest, "project name is required")
		case domain.ErrProjectAlreadyExists:
			respondWithError(w, http.StatusConflict, "project with this name already exists")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to create project")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, project)
}

// UpdateProject handles PATCH /api/projects/{id}
// @Summary Update a project
// @Description Update an existing project's name
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Param project body object true "Project update request" schema="{name:string}"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id} [put]
func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing project ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID format")
		return
	}

	var request struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	project, err := h.projectService.UpdateProject(ctx, id, request.Name)
	if err != nil {
		switch err {
		case domain.ErrProjectNotFound:
			respondWithError(w, http.StatusNotFound, "project not found")
		case domain.ErrInvalidProjectName:
			respondWithError(w, http.StatusBadRequest, "project name is required")
		case domain.ErrProjectAlreadyExists:
			respondWithError(w, http.StatusConflict, "project with this name already exists")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to update project")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, project)
}

// DeleteProject handles DELETE /api/projects/{id}
// @Summary Delete a project
// @Description Delete a project by its ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path int true "Project ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{id} [delete]
func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing project ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID format")
		return
	}

	ctx := r.Context()
	err = h.projectService.DeleteProject(ctx, id)
	if err != nil {
		switch err {
		case domain.ErrProjectNotFound:
			respondWithError(w, http.StatusNotFound, "project not found")
		case domain.ErrInvalidProjectName:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to delete project")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "project deleted successfully"})
}

// Helper functions for HTTP responses
func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Log error but can't write to response as headers already sent
		fmt.Printf("Error encoding JSON response: %v\n", err)
	}
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}
