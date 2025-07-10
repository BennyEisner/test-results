package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/failure/domain/ports"
)

// FailureHandler handles HTTP requests for failures
type FailureHandler struct {
	Service ports.FailureService
}

// NewFailureHandler creates a new FailureHandler
func NewFailureHandler(service ports.FailureService) *FailureHandler {
	return &FailureHandler{Service: service}
}

// GetFailureByID handles GET /failures/{id}
// @Summary Get failure by ID
// @Description Retrieve a failure by its unique identifier
// @Tags failures
// @Accept json
// @Produce json
// @Param id path int true "Failure ID"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /failures/{id} [get]
func (h *FailureHandler) GetFailureByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid failure ID")
		return
	}

	ctx := r.Context()
	failure, err := h.Service.GetFailure(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, failure)
}

// GetFailureByExecution handles GET /executions/{executionID}/failure
// @Summary Get failure by execution ID
// @Description Retrieve a failure associated with a specific test execution
// @Tags failures
// @Accept json
// @Produce json
// @Param executionID path int true "Execution ID"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /executions/{executionID}/failure [get]
func (h *FailureHandler) GetFailureByExecution(w http.ResponseWriter, r *http.Request) {
	executionIDStr := r.PathValue("executionID")
	executionID, err := strconv.ParseInt(executionIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid execution ID")
		return
	}

	ctx := r.Context()
	failure, err := h.Service.GetFailureByExecution(ctx, executionID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, failure)
}

// CreateFailure handles POST /executions/{executionID}/failures
// @Summary Create a new failure
// @Description Create a new failure for a specific test execution
// @Tags failures
// @Accept json
// @Produce json
// @Param executionID path int true "Execution ID"
// @Param failure body object true "Failure creation request" schema="{message:string,type:string,details:string}"
// @Success 201 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /executions/{executionID}/failures [post]
func (h *FailureHandler) CreateFailure(w http.ResponseWriter, r *http.Request) {
	executionIDStr := r.PathValue("executionID")
	executionID, err := strconv.ParseInt(executionIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid execution ID")
		return
	}

	var input struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Details string `json:"details"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	failure, err := h.Service.CreateFailure(ctx, executionID, input.Message, input.Type, input.Details)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, failure)
}

// UpdateFailure handles PUT /failures/{id}
// @Summary Update a failure
// @Description Update an existing failure's details
// @Tags failures
// @Accept json
// @Produce json
// @Param id path int true "Failure ID"
// @Param failure body object true "Failure update request" schema="{message:string,type:string,details:string}"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /failures/{id} [put]
func (h *FailureHandler) UpdateFailure(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid failure ID")
		return
	}

	var input struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Details string `json:"details"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	updatedFailure, err := h.Service.UpdateFailure(ctx, id, input.Message, input.Type, input.Details)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, updatedFailure)
}

// DeleteFailure handles DELETE /failures/{id}
// @Summary Delete a failure
// @Description Delete a failure by its ID
// @Tags failures
// @Accept json
// @Produce json
// @Param id path int true "Failure ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /failures/{id} [delete]
func (h *FailureHandler) DeleteFailure(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid failure ID")
		return
	}

	ctx := r.Context()
	if err := h.Service.DeleteFailure(ctx, id); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helper functions for HTTP responses
func respondWithError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
