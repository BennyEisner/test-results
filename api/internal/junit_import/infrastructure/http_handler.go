package http

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// JUnitImportHandler handles HTTP requests for JUnit imports
type JUnitImportHandler struct {
	Service interface{}
}

// NewJUnitImportHandler creates a new JUnitImportHandler
func NewJUnitImportHandler(service interface{}) *JUnitImportHandler {
	return &JUnitImportHandler{Service: service}
}

// ProcessJUnitData handles POST /import/junit
// @Summary Import JUnit test data
// @Description Process and import JUnit XML test results for a specific project and test suite
// @Tags junit-import
// @Accept json
// @Produce json
// @Param project_id query int true "Project ID"
// @Param suite_id query int true "Test Suite ID"
// @Param junitData body object true "JUnit test data"
// @Success 201 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /import/junit [post]
func (h *JUnitImportHandler) ProcessJUnitData(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	_, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
		return
	}

	suiteIDStr := r.URL.Query().Get("suite_id")
	_, err = strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid suite_id", http.StatusBadRequest)
		return
	}

	var junitData interface{}
	if err := json.NewDecoder(r.Body).Decode(&junitData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// Temporarily disabled
	http.Error(w, "JUnit import service temporarily disabled", http.StatusServiceUnavailable)
}
