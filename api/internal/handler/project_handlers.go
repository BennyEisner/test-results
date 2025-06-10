package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BennyEisner/test-results/internal/utils"
)

// HandleDBTest handles database test endpoint
func HandleDBTest(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Database error: %v", err)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Database connection successful. Projects count: %d", count)
}

// HandleProjects handles project collection endpoints
func HandleProjects(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodGet:
		getProjects(w, r, db)
	case http.MethodPost:
		createProject(w, r, db)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleProjectByPath handles operations on specific projects
func HandleProjectByPath(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Special case for create endpoint
	if r.URL.Path == "/api/projects/create" && r.Method == http.MethodPost {
		createProject(w, r, db)
		return
	}

	// Extract ID from the path
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/projects/"), "/")
	if len(pathSegments) != 1 || pathSegments[0] == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID in URL")
		return
	}

	idStr := pathSegments[0]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format")
		return
	}

	switch r.Method {
	case http.MethodGet:
		getProjectByID(w, r, id, db)
	case http.MethodPatch:
		updateProject(w, r, id, db)
	case http.MethodDelete:
		deleteProject(w, r, id, db)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// Get a project by ID
func getProjectByID(w http.ResponseWriter, r *http.Request, id int, db *sql.DB) {
	var p utils.Project
	err := db.QueryRow("SELECT id, name FROM projects WHERE id = $1", id).Scan(&p.ID, &p.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Project not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, p)
}

// Get all projects
func getProjects(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT id, name FROM projects ORDER BY id")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	defer rows.Close()

	projects := []utils.Project{}
	for rows.Next() {
		var p utils.Project
		if err := rows.Scan(&p.ID, &p.Name); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Scan error: "+err.Error())
			return
		}
		projects = append(projects, p)
	}
	utils.RespondWithJSON(w, http.StatusOK, projects)
}

// Create a new project
func createProject(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed")
		return
	}

	var p utils.Project
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	defer r.Body.Close()

	if p.Name == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Project name is required")
		return
	}

	var id int
	err := db.QueryRow("INSERT INTO projects(name) VALUES($1) RETURNING id", p.Name).Scan(&id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}
	p.ID = id
	utils.RespondWithJSON(w, http.StatusCreated, p)
}

// Delete a project by ID
func deleteProject(w http.ResponseWriter, r *http.Request, id int, db *sql.DB) {
	result, err := db.Exec("DELETE FROM projects WHERE id = $1", id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error checking delete result: "+err.Error())
		return
	}

	if rowsAffected == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Project deleted successfully"})
}

// Update a project by ID
func updateProject(w http.ResponseWriter, r *http.Request, id int, db *sql.DB) {
	var updateData map[string]interface{}
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&updateData); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}
	defer r.Body.Close()

	if len(updateData) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "No fields provided for update")
		return
	}

	// Check if project exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
		return
	}

	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Project not found")
		return
	}

	// SQL Update from provided parameters
	updateFields := []string{}
	values := []interface{}{}
	valueIndex := 1

	if name, ok := updateData["name"].(string); ok {
		if name == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Project name cannot be empty")
			return
		}
		updateFields = append(updateFields, fmt.Sprintf("name = $%d", valueIndex))
		values = append(values, name)
		valueIndex++
	}

	if len(updateFields) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "No valid fields provided for update")
		return
	}

	query := fmt.Sprintf("UPDATE projects SET %s WHERE id=$%d RETURNING id, name",
		strings.Join(updateFields, ", "), valueIndex)
	values = append(values, id)

	var updatedProject utils.Project
	err = db.QueryRow(query, values...).Scan(&updatedProject.ID, &updatedProject.Name)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Update failed: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, updatedProject)
}
