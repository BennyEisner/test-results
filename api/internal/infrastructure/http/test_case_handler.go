package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain"
)

type TestCaseHandler struct {
	service domain.TestCaseService
}

func NewTestCaseHandler(service domain.TestCaseService) *TestCaseHandler {
	return &TestCaseHandler{service: service}
}

func (h *TestCaseHandler) GetTestCaseByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing test case ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid test case ID format")
		return
	}

	testCase, err := h.service.GetTestCaseByID(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrTestCaseNotFound:
			respondWithError(w, http.StatusNotFound, "test case not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, testCase)
}

func (h *TestCaseHandler) GetTestCasesBySuiteID(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.PathValue("suiteID")
	if suiteIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing suite ID")
		return
	}

	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid suite ID format")
		return
	}

	testCases, err := h.service.GetTestCasesBySuiteID(r.Context(), suiteID)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, testCases)
}

func (h *TestCaseHandler) GetTestCaseByName(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.PathValue("suiteID")
	if suiteIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing suite ID")
		return
	}

	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid suite ID format")
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		respondWithError(w, http.StatusBadRequest, "missing name parameter")
		return
	}

	testCase, err := h.service.GetTestCaseByName(r.Context(), suiteID, name)
	if err != nil {
		switch err {
		case domain.ErrTestCaseNotFound:
			respondWithError(w, http.StatusNotFound, "test case not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, testCase)
}

func (h *TestCaseHandler) CreateTestCase(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Name      string `json:"name"`
		Classname string `json:"classname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	suiteIDStr := r.PathValue("suiteID")
	if suiteIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing suite ID")
		return
	}

	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid suite ID format")
		return
	}

	testCase, err := h.service.CreateTestCase(r.Context(), suiteID, request.Name, request.Classname)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "name and classname are required")
		case domain.ErrDuplicateTestCase:
			respondWithError(w, http.StatusConflict, "test case with this name already exists in this suite")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to create test case")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, testCase)
}

func (h *TestCaseHandler) UpdateTestCase(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing test case ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid test case ID format")
		return
	}

	var request struct {
		Name      string `json:"name"`
		Classname string `json:"classname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	testCase, err := h.service.UpdateTestCase(r.Context(), id, request.Name, request.Classname)
	if err != nil {
		switch err {
		case domain.ErrTestCaseNotFound:
			respondWithError(w, http.StatusNotFound, "test case not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "name and classname are required")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to update test case")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, testCase)
}

func (h *TestCaseHandler) DeleteTestCase(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing test case ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid test case ID format")
		return
	}

	if err := h.service.DeleteTestCase(r.Context(), id); err != nil {
		switch err {
		case domain.ErrTestCaseNotFound:
			respondWithError(w, http.StatusNotFound, "test case not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to delete test case")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "test case deleted successfully"})
}
