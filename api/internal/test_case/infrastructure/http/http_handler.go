package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/test_case/domain/ports"
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
// @Summary Get test case by ID
// @Description Retrieve a test case by its unique identifier
// @Tags test-cases
// @Accept json
// @Produce json
// @Param id query int true "Test Case ID"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /test-cases [get]
func (h *TestCaseHandler) GetTestCaseByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	testCase, err := h.Service.GetTestCase(r.Context(), id)
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
// @Summary Get test cases by suite ID
// @Description Retrieve all test cases for a specific test suite
// @Tags test-cases
// @Accept json
// @Produce json
// @Param suite_id query int true "Test Suite ID"
// @Success 200 {array} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /test-cases [get]
func (h *TestCaseHandler) GetTestCasesBySuiteID(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.URL.Query().Get("suite_id")
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid suite_id", http.StatusBadRequest)
		return
	}
	testCases, err := h.Service.GetTestCasesBySuite(r.Context(), suiteID)
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
// @Summary Create a new test case
// @Description Create a new test case for a test suite
// @Tags test-cases
// @Accept json
// @Produce json
// @Param testCase body object true "Test case creation request" schema="{suite_id:int,name:string,classname:string}"
// @Success 201 {object} models.TestCase
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /test-cases [post]
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
// @Summary Update a test case
// @Description Update an existing test case's name and classname
// @Tags test-cases
// @Accept json
// @Produce json
// @Param id query int true "Test Case ID"
// @Param testCase body object true "Test case update request" schema="{name:string,classname:string}"
// @Success 200 {object} models.TestCase
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /test-cases [put]
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
// @Summary Delete a test case
// @Description Delete a test case by its ID
// @Tags test-cases
// @Accept json
// @Produce json
// @Param id query int true "Test Case ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /test-cases [delete]
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
