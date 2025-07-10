package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/build_test_case_execution/domain/models"
	"github.com/BennyEisner/test-results/internal/build_test_case_execution/domain/ports"
)

// BuildTestCaseExecutionHandler handles HTTP requests for build test case executions
type BuildTestCaseExecutionHandler struct {
	Service ports.BuildTestCaseExecutionService
}

// NewBuildTestCaseExecutionHandler creates a new handler
func NewBuildTestCaseExecutionHandler(service ports.BuildTestCaseExecutionService) *BuildTestCaseExecutionHandler {
	return &BuildTestCaseExecutionHandler{Service: service}
}

// GetExecutionByID handles GET /executions/{id}
// @Summary Get execution by ID
// @Description Retrieve a build test case execution by its unique identifier
// @Tags executions
// @Accept json
// @Produce json
// @Param id path int true "Execution ID"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /executions/{id} [get]
func (h *BuildTestCaseExecutionHandler) GetExecutionByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid execution ID")
		return
	}

	ctx := r.Context()
	execution, err := h.Service.GetExecutionByID(ctx, id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, execution)
}

// GetExecutionsByBuildID handles GET /builds/{buildID}/executions
// @Summary Get executions by build ID
// @Description Retrieve all test case executions for a specific build
// @Tags executions
// @Accept json
// @Produce json
// @Param buildID path int true "Build ID"
// @Success 200 {array} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /builds/{buildID}/executions [get]
func (h *BuildTestCaseExecutionHandler) GetExecutionsByBuildID(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.PathValue("buildID")
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid build ID")
		return
	}

	ctx := r.Context()
	executions, err := h.Service.GetExecutionsByBuildID(ctx, buildID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, executions)
}

// CreateExecution handles POST /builds/{buildID}/executions
// @Summary Create a new execution
// @Description Create a new test case execution for a specific build
// @Tags executions
// @Accept json
// @Produce json
// @Param buildID path int true "Build ID"
// @Param execution body object true "Execution creation request"
// @Success 201 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /builds/{buildID}/executions [post]
func (h *BuildTestCaseExecutionHandler) CreateExecution(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.PathValue("buildID")
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid build ID")
		return
	}

	var input models.BuildExecutionInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	execution, err := h.Service.CreateExecution(ctx, buildID, &input)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, execution)
}

// UpdateExecution handles PUT /executions/{id}
// @Summary Update an execution
// @Description Update an existing test case execution
// @Tags executions
// @Accept json
// @Produce json
// @Param id path int true "Execution ID"
// @Param execution body object true "Execution update request"
// @Success 200 {object} object
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /executions/{id} [put]
func (h *BuildTestCaseExecutionHandler) UpdateExecution(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid execution ID")
		return
	}

	var execution models.BuildTestCaseExecution
	if err := json.NewDecoder(r.Body).Decode(&execution); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ctx := r.Context()
	updatedExecution, err := h.Service.UpdateExecution(ctx, id, &execution)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, updatedExecution)
}

// DeleteExecution handles DELETE /executions/{id}
// @Summary Delete an execution
// @Description Delete a test case execution by its ID
// @Tags executions
// @Accept json
// @Produce json
// @Param id path int true "Execution ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /executions/{id} [delete]
func (h *BuildTestCaseExecutionHandler) DeleteExecution(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid execution ID")
		return
	}

	ctx := r.Context()
	if err := h.Service.DeleteExecution(ctx, id); err != nil {
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
