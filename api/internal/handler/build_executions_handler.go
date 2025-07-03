package handler

import (
	"net/http"
	"strconv"

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

func (beh *BuildExecutionHandler) GetBuildExecutions(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.PathValue("id")
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID format: "+err.Error())
		return
	}

	executions, err := beh.service.GetBuildExecutions(buildID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching build executions: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, executions)
}
