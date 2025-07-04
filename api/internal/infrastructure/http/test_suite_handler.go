package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain/ports"
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
func (h *TestSuiteHandler) GetTestSuiteByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	suite, err := h.Service.GetTestSuiteByID(r.Context(), id)
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
func (h *TestSuiteHandler) GetTestSuitesByProjectID(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("project_id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid project_id", http.StatusBadRequest)
		return
	}
	suites, err := h.Service.GetTestSuitesByProjectID(r.Context(), projectID)
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
	suite, err := h.Service.CreateTestSuite(r.Context(), req.ProjectID, req.Name, req.ParentID)
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
