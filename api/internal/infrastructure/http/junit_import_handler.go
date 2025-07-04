package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// JUnitImportHandler handles HTTP requests for JUnit imports
type JUnitImportHandler struct {
	Service ports.JUnitImportService
}

// NewJUnitImportHandler creates a new JUnitImportHandler
func NewJUnitImportHandler(service ports.JUnitImportService) *JUnitImportHandler {
	return &JUnitImportHandler{Service: service}
}

// ProcessJUnitData handles POST /import/junit
func (h *JUnitImportHandler) ProcessJUnitData(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
		return
	}

	suiteIDStr := r.URL.Query().Get("suite_id")
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid suite_id", http.StatusBadRequest)
		return
	}

	var junitData models.JUnitTestSuites
	if err := json.NewDecoder(r.Body).Decode(&junitData); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	build, err := h.Service.ProcessJUnitData(r.Context(), projectID, suiteID, &junitData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(build); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
