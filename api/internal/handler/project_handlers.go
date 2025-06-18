package handler

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/BennyEisner/test-results/internal/service"
	"github.com/BennyEisner/test-results/internal/utils"
)

// ProjectHandler holds the project service.
type ProjectHandler struct {
	service service.ProjectServiceInterface
}

// NewProjectHandler creates a new ProjectHandler.
func NewProjectHandler(s service.ProjectServiceInterface) *ProjectHandler {
	return &ProjectHandler{service: s}
}

// Helper to convert models.Project to utils.Project DTO
func toAPIProject(m *models.Project) utils.Project {
	return utils.Project{
		ID:   int(m.ID), // Convert int64 to int
		Name: m.Name,
	}
}

// Helper to convert slice of models.Project to slice of utils.Project DTO
func toAPIProjects(ms []models.Project) []utils.Project {
	apiProjects := make([]utils.Project, len(ms))
	for i, m := range ms {
		apiProjects[i] = toAPIProject(&m)
	}
	return apiProjects
}

// HandleDBTest handles database test endpoint
func (ph *ProjectHandler) HandleDBTest(w http.ResponseWriter, r *http.Request) {
	count, err := ph.service.GetDBTestProjectCount()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Database error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Database connection successful. Projects count: %d", count)
}

// HandleProjects handles project collection endpoints
func (ph *ProjectHandler) HandleProjects(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers to allow frontend access
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS") // Add other methods if needed for this specific path
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight OPTIONS requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case http.MethodGet:
		ph.getProjects(w, r)
	case http.MethodPost:
		ph.createProject(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleProjectByPath handles operations on specific projects
func (ph *ProjectHandler) HandleProjectByPath(w http.ResponseWriter, r *http.Request) {
	// Special case for create endpoint - this logic might be better handled by distinct routing
	// if r.URL.Path == "/api/projects/create" && r.Method == http.MethodPost {
	// 	ph.createProject(w, r) // createProject is now part of HandleProjects
	// 	return
	// }

	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/projects/"), "/")
	if len(pathSegments) != 1 || pathSegments[0] == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID in URL")
		return
	}

	idStr := pathSegments[0]
	// Project ID in models.Project is int64, but utils.Project and current path parsing use int.
	// For consistency with service layer, parse to int64.
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	switch r.Method {
	case http.MethodGet:
		ph.getProjectByID(w, r, id)
	case http.MethodPatch: // Assuming PATCH for updates
		ph.updateProject(w, r, id)
	case http.MethodDelete:
		ph.deleteProject(w, r, id)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (ph *ProjectHandler) getProjectByID(w http.ResponseWriter, r *http.Request, id int64) {
	projectModel, err := ph.service.GetProjectByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Project not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, toAPIProject(projectModel))
}

func (ph *ProjectHandler) getProjects(w http.ResponseWriter, r *http.Request) {
	log.Println("✅ getProjects called via ProjectHandler") // Keep log for now
	projectModels, err := ph.service.GetAllProjects()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, toAPIProjects(projectModels))
}

func (ph *ProjectHandler) createProject(w http.ResponseWriter, r *http.Request) {
	// POST is already checked by HandleProjects method or router
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var projectName string
	contentType := r.Header.Get("Content-Type")

	if strings.Contains(contentType, "application/json") {
		var payload struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON request: "+err.Error())
			return
		}
		projectName = payload.Name
	} else if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		var projectXML struct {
			XMLName xml.Name `xml:"project"`
			Name    string   `xml:"name"`
		}
		if err := xml.Unmarshal(body, &projectXML); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid XML request: "+err.Error())
			return
		}
		projectName = projectXML.Name
	} else {
		utils.RespondWithError(w, http.StatusUnsupportedMediaType, "Content type must be application/xml or application/json")
		return
	}

	if strings.TrimSpace(projectName) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Project Name is required")
		return
	}

	createdProjectModel, err := ph.service.CreateProject(projectName)
	if err != nil {
		// Consider more specific error mapping if service returns typed errors
		utils.RespondWithError(w, http.StatusInternalServerError, "Database Error: "+err.Error())
		return
	}

	apiProject := toAPIProject(createdProjectModel)

	acceptHeader := r.Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/xml") {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusCreated)
		// Using models.ProjectXML for consistency if it's defined for XML marshalling
		xmlResponse := models.ProjectXML{Project: models.Project{ID: createdProjectModel.ID, Name: createdProjectModel.Name}}
		xml.NewEncoder(w).Encode(xmlResponse)
	} else {
		utils.RespondWithJSON(w, http.StatusCreated, apiProject)
	}
}

func (ph *ProjectHandler) deleteProject(w http.ResponseWriter, r *http.Request, id int64) {
	rowsAffected, err := ph.service.DeleteProject(id)
	if err != nil {
		if err == sql.ErrNoRows { // Service might not return this directly, depends on its impl.
			utils.RespondWithError(w, http.StatusNotFound, "Project not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		}
		return
	}
	if rowsAffected == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Project not found or already deleted")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Project deleted successfully"})
}

func (ph *ProjectHandler) updateProject(w http.ResponseWriter, r *http.Request, id int64) {
	// Assuming PATCH, body parsing for name
	var payload struct {
		Name string `json:"name" xml:"name"` // Support both JSON and XML for update
	}

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid JSON request: "+err.Error())
			return
		}
	} else if strings.Contains(contentType, "application/xml") {
		if err := xml.NewDecoder(r.Body).Decode(&payload); err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid XML request: "+err.Error())
			return
		}
	} else {
		utils.RespondWithError(w, http.StatusUnsupportedMediaType, "Content type must be application/xml or application/json for update")
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(payload.Name) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Project Name is required for update")
		return
	}

	updatedProjectModel, err := ph.service.UpdateProject(id, payload.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Project not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Update failed: "+err.Error())
		}
		return
	}

	apiProject := toAPIProject(updatedProjectModel)
	utils.RespondWithJSON(w, http.StatusOK, apiProject)
}
