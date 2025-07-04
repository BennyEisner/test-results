package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// TestCaseHandler handles HTTP requests for test cases
type TestCaseHandler struct {
	Service ports.TestCaseService
}

// NewTestCaseHandler creates a new TestCaseHandler
func NewTestCaseHandler(service ports.TestCaseService) *TestCaseHandler {
	return &TestCaseHandler{Service: service}
}

// GetTestCaseByID handles GET /test-cases/{id}
func (h *TestCaseHandler) GetTestCaseByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	testCase, err := h.Service.GetTestCaseByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if testCase == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(testCase); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// GetTestCasesBySuiteID handles GET /test-suites/{suiteID}/test-cases
func (h *TestCaseHandler) GetTestCasesBySuiteID(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.URL.Query().Get("suite_id")
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid suite_id", http.StatusBadRequest)
		return
	}
	testCases, err := h.Service.GetTestCasesBySuiteID(r.Context(), suiteID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := json.NewEncoder(w).Encode(testCases); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// CreateTestCase handles POST /test-cases
func (h *TestCaseHandler) CreateTestCase(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SuiteID   int64  `json:"suite_id"`
		Name      string `json:"name"`
		Classname string `json:"classname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	testCase, err := h.Service.CreateTestCase(r.Context(), req.SuiteID, req.Name, req.Classname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(testCase); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// UpdateTestCase handles PUT /test-cases/{id}
func (h *TestCaseHandler) UpdateTestCase(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	var req struct {
		Name      string `json:"name"`
		Classname string `json:"classname"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	testCase, err := h.Service.UpdateTestCase(r.Context(), id, req.Name, req.Classname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if testCase == nil {
		http.NotFound(w, r)
		return
	}
	if err := json.NewEncoder(w).Encode(testCase); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}

// DeleteTestCase handles DELETE /test-cases/{id}
func (h *TestCaseHandler) DeleteTestCase(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteTestCase(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
