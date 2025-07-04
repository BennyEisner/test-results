package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain"
)

type BuildTestCaseExecutionHandler struct {
	svc domain.BuildTestCaseExecutionService
}

func NewBuildTestCaseExecutionHandler(svc domain.BuildTestCaseExecutionService) *BuildTestCaseExecutionHandler {
	return &BuildTestCaseExecutionHandler{svc: svc}
}

func (h *BuildTestCaseExecutionHandler) GetExecutionByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid execution ID")
		return
	}
	ctx := r.Context()
	execution, err := h.svc.GetExecutionByID(ctx, id)
	if err != nil {
		switch err {
		case domain.ErrExecutionNotFound:
			respondWithError(w, http.StatusNotFound, "execution not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal error")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, execution)
}

func (h *BuildTestCaseExecutionHandler) GetExecutionsByBuildID(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.PathValue("id")
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid build ID")
		return
	}
	ctx := r.Context()
	executions, err := h.svc.GetExecutionsByBuildID(ctx, buildID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to get executions")
		return
	}
	respondWithJSON(w, http.StatusOK, executions)
}

func (h *BuildTestCaseExecutionHandler) CreateExecution(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.PathValue("id")
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid build ID")
		return
	}
	var input domain.BuildExecutionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx := r.Context()
	execution, err := h.svc.CreateExecution(ctx, buildID, &input)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to create execution")
		}
		return
	}
	respondWithJSON(w, http.StatusCreated, execution)
}

func (h *BuildTestCaseExecutionHandler) UpdateExecution(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid execution ID")
		return
	}
	var execution domain.BuildTestCaseExecution
	if err := json.NewDecoder(r.Body).Decode(&execution); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	ctx := r.Context()
	updatedExecution, err := h.svc.UpdateExecution(ctx, id, &execution)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to update execution")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, updatedExecution)
}

func (h *BuildTestCaseExecutionHandler) DeleteExecution(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid execution ID")
		return
	}
	ctx := r.Context()
	if err := h.svc.DeleteExecution(ctx, id); err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to delete execution")
		}
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "execution deleted successfully"})
}
