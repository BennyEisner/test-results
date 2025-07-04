package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain"
)

type FailureHandler struct {
	service domain.FailureService
}

func NewFailureHandler(service domain.FailureService) *FailureHandler {
	return &FailureHandler{service: service}
}

func (h *FailureHandler) GetFailureByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing failure ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid failure ID format")
		return
	}

	failure, err := h.service.GetFailureByID(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrFailureNotFound:
			respondWithError(w, http.StatusNotFound, "failure not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, failure)
}

func (h *FailureHandler) GetFailureByExecutionID(w http.ResponseWriter, r *http.Request) {
	executionIDStr := r.PathValue("executionID")
	if executionIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing execution ID")
		return
	}

	executionID, err := strconv.ParseInt(executionIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid execution ID format")
		return
	}

	failure, err := h.service.GetFailureByExecutionID(r.Context(), executionID)
	if err != nil {
		switch err {
		case domain.ErrFailureNotFound:
			respondWithError(w, http.StatusNotFound, "failure not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, failure)
}

func (h *FailureHandler) CreateFailure(w http.ResponseWriter, r *http.Request) {
	var failure domain.Failure
	if err := json.NewDecoder(r.Body).Decode(&failure); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createdFailure, err := h.service.CreateFailure(r.Context(), &failure)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to create failure")
		}
		return
	}

	respondWithJSON(w, http.StatusCreated, createdFailure)
}

func (h *FailureHandler) UpdateFailure(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing failure ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid failure ID format")
		return
	}

	var failure domain.Failure
	if err := json.NewDecoder(r.Body).Decode(&failure); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedFailure, err := h.service.UpdateFailure(r.Context(), id, &failure)
	if err != nil {
		switch err {
		case domain.ErrFailureNotFound:
			respondWithError(w, http.StatusNotFound, "failure not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to update failure")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, updatedFailure)
}

func (h *FailureHandler) DeleteFailure(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		respondWithError(w, http.StatusBadRequest, "missing failure ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid failure ID format")
		return
	}

	if err := h.service.DeleteFailure(r.Context(), id); err != nil {
		switch err {
		case domain.ErrFailureNotFound:
			respondWithError(w, http.StatusNotFound, "failure not found")
		case domain.ErrInvalidInput:
			respondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			respondWithError(w, http.StatusInternalServerError, "failed to delete failure")
		}
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "failure deleted successfully"})
}
