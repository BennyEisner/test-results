package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain"
)

type TestSuiteHandler struct {
	svc domain.TestSuiteService
}

func NewTestSuiteHandler(svc domain.TestSuiteService) *TestSuiteHandler {
	return &TestSuiteHandler{svc: svc}
}

func (h *TestSuiteHandler) GetTestSuiteByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid test suite ID")
		return
	}
	ctx := r.Context()
	ts, err := h.svc.GetTestSuiteByID(ctx, id)
	if err != nil {
		switch err {
		case domain.ErrTestSuiteNotFound:
			respondWithError(w, http.StatusNotFound, "test suite not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, ts)
}

func (h *TestSuiteHandler) GetTestSuitesByProjectID(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}
	ctx := r.Context()
	ts, err := h.svc.GetTestSuitesByProjectID(ctx, projectID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get test suites")
		return
	}
	respondWithJSON(w, http.StatusOK, ts)
}

func (h *TestSuiteHandler) CreateTestSuite(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid project ID")
		return
	}
	var req struct {
		Name     string `json:"name"`
		ParentID *int64 `json:"parent_id,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx := r.Context()
	ts, err := h.svc.CreateTestSuite(ctx, projectID, req.Name, req.ParentID)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		case domain.ErrDuplicateTestSuite:
			respondWithError(w, http.StatusConflict, "test suite already exists")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to create test suite")
		}
		return
	}
	respondWithJSON(w, http.StatusCreated, ts)
}

func (h *TestSuiteHandler) UpdateTestSuite(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid test suite ID")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx := r.Context()
	ts, err := h.svc.UpdateTestSuite(ctx, id, req.Name)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to update test suite")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, ts)
}

func (h *TestSuiteHandler) DeleteTestSuite(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid test suite ID")
		return
	}
	ctx := r.Context()
	if err := h.svc.DeleteTestSuite(ctx, id); err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to delete test suite")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "test suite deleted successfully"})
}
