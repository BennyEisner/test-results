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

func (ph *ProjectHandler) GetProjectByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format")
		return
	}

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

func (ph *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name != "" {
		projectModel, err := ph.service.GetProjectByName(name)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.RespondWithError(w, http.StatusNotFound, "Project not found")
			} else {
				utils.RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
			}
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, toAPIProject(projectModel))
		return
	}
	log.Println("✅ GetProjects called via ProjectHandler") // Keep log for now
	projectModels, err := ph.service.GetAllProjects()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, toAPIProjects(projectModels))
}

func (ph *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
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
		utils.RespondWithError(w, http.StatusInternalServerError, "Database Error: "+err.Error())
		return
	}

	apiProject := toAPIProject(createdProjectModel)

	acceptHeader := r.Header.Get("Accept")
	if strings.Contains(acceptHeader, "application/xml") {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusCreated)
		xmlResponse := models.ProjectXML{Project: models.Project{ID: createdProjectModel.ID, Name: createdProjectModel.Name}}
		if err := xml.NewEncoder(w).Encode(xmlResponse); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	} else {
		utils.RespondWithJSON(w, http.StatusCreated, apiProject)
	}
}

func (ph *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	rowsAffected, err := ph.service.DeleteProject(id)
	if err != nil {
		if err == sql.ErrNoRows {
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

func (ph *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	var payload struct {
		Name string `json:"name" xml:"name"`
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
