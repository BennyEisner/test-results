package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/BennyEisner/test-results/internal/service" // Import the service package
	"github.com/BennyEisner/test-results/internal/utils"
)

// BuildExecutionHandler holds the build execution service.
type BuildExecutionHandler struct {
	service service.BuildExecutionServiceInterface
}

// NewBuildExecutionHandler creates a new BuildExecutionHandler.
func NewBuildExecutionHandler(s service.BuildExecutionServiceInterface) *BuildExecutionHandler {
	return &BuildExecutionHandler{service: s}
}

// HandleBuildExecutions handles GET and POST requests for test case executions of a specific build.
// Expected path: /api/builds/{buildId}/executions
func (beh *BuildExecutionHandler) HandleBuildExecutions(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/builds/"), "/")

	if len(pathSegments) < 2 || pathSegments[0] == "" || pathSegments[1] != "executions" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL path for build executions. Expected /api/builds/{buildId}/executions")
		return
	}

	buildIDStr := pathSegments[0]
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID format: "+err.Error())
		return
	}

	buildExists, err := beh.service.CheckBuildExists(buildID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error checking build existence: "+err.Error())
		return
	}
	if !buildExists {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Build with ID %d not found", buildID))
		return
	}

	switch r.Method {
	case http.MethodGet:
		beh.getBuildExecutions(w, r, buildID)
	case http.MethodPost:
		beh.createBuildExecutions(w, r, buildID)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed for this endpoint")
	}
}

func (beh *BuildExecutionHandler) getBuildExecutions(w http.ResponseWriter, r *http.Request, buildID int64) {
	executions, err := beh.service.GetBuildExecutions(buildID)
	if err != nil {
		// The service's GetBuildExecutions already returns []models.BuildExecutionDetail (formerly handler.BuildExecutionDetail)
		// So, no specific error handling for sql.ErrNoRows here, service handles it.
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching build executions: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, executions)
}

func (beh *BuildExecutionHandler) createBuildExecutions(w http.ResponseWriter, r *http.Request, buildID int64) {
	var inputs []models.BuildExecutionInput // Use DTO from models package
	if err := json.NewDecoder(r.Body).Decode(&inputs); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	if len(inputs) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Request must contain at least one execution result.")
		return
	}

	// Validation of inputs (e.g., status values) can remain here or be moved more into the service
	// For now, keeping basic input format validation here. Service handles DB existence checks.
	validStatuses := map[string]bool{"passed": true, "failed": true, "skipped": true, "error": true}
	for _, input := range inputs {
		statusToCheck := strings.ToLower(input.Status)
		if !validStatuses[statusToCheck] {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid status '%s' for test case ID %d. Must be one of: passed, failed, skipped, error.", input.Status, input.TestCaseID))
			return // Fail fast on first invalid status in input
		}
	}

	createdExecutionModels, processingErrors, err := beh.service.CreateBuildExecutions(buildID, inputs)
	if err != nil {
		// This is for catastrophic errors like DB transaction start failure
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process build executions: "+err.Error())
		return
	}

	// Construct response based on service output
	// The service returns models.BuildTestCaseExecution, which is simpler than the full detail.
	// For POST response, a summary is often sufficient.
	response := map[string]interface{}{
		"message":              fmt.Sprintf("%d execution results processed.", len(inputs)),
		"successful_creations": len(createdExecutionModels), // Count of successfully created models
		"errors_encountered":   len(processingErrors),
	}

	if len(processingErrors) > 0 {
		response["error_details"] = processingErrors
		utils.RespondWithJSON(w, http.StatusMultiStatus, response)
	} else {
		utils.RespondWithJSON(w, http.StatusCreated, response)
	}
}
