package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/test_suite/domain/ports"
)

// TestSuiteHandler handles HTTP requests for test suites
type TestSuiteHandler struct {
	Service ports.TestSuiteService
}

// NewTestSuiteHandler creates a new TestSuiteHandler
func NewTestSuiteHandler(service ports.TestSuiteService) *TestSuiteHandler {
	return &TestSuiteHandler{Service: service}
}

// GetTestSuiteByID handles GET /test-suites/{id}
// @Summary Get test suite by ID
// @Description Retrieve a test suite by its unique identifier
// @Tags test-suites
// @Accept json
// @Produce json
// @Param id query int true "Test Suite ID"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /test-suites [get]
func (h *TestSuiteHandler) GetTestSuiteByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	suite, err := h.Service.GetTestSuite(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if suite == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(suite); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetTestSuitesByProjectID handles GET /projects/{projectID}/test-suites
// @Summary Get test suites by project ID
// @Description Retrieve all test suites for a specific project
// @Tags test-suites
// @Accept json
// @Produce json
// @Param project_id query int true "Project ID"
// @Success 200 {array} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /test-suites [get]
func (h *TestSuiteHandler) GetTestSuitesByProjectID(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
		return
	}
	suites, err := h.Service.GetTestSuitesByProject(r.Context(), projectID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(suites); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateTestSuite handles POST /test-suites
// @Summary Create a new test suite
// @Description Create a new test suite for a project
// @Tags test-suites
// @Accept json
// @Produce json
// @Param suite body object true "Test suite creation request" schema="{project_id:int,name:string,parent_id:int}"
// @Success 201 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /test-suites [post]
func (h *TestSuiteHandler) CreateTestSuite(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ProjectID int64  `json:"project_id"`
		Name      string `json:"name"`
		ParentID  *int64 `json:"parent_id,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	suite, err := h.Service.CreateTestSuite(r.Context(), req.ProjectID, req.Name, req.ParentID, 0.0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(suite); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateTestSuite handles PUT /test-suites/{id}
// @Summary Update a test suite
// @Description Update an existing test suite's name
// @Tags test-suites
// @Accept json
// @Produce json
// @Param id query int true "Test Suite ID"
// @Param suite body object true "Test suite update request" schema="{name:string}"
// @Success 200 {object} models.TestSuite
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /test-suites [put]
func (h *TestSuiteHandler) UpdateTestSuite(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	suite, err := h.Service.UpdateTestSuite(r.Context(), id, req.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if suite == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(suite); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteTestSuite handles DELETE /test-suites/{id}
// @Summary Delete a test suite
// @Description Delete a test suite by its ID
// @Tags test-suites
// @Accept json
// @Produce json
// @Param id query int true "Test Suite ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /test-suites [delete]
func (h *TestSuiteHandler) DeleteTestSuite(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteTestSuite(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
