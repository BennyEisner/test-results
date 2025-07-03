package handler

import (
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/service"
	"github.com/BennyEisner/test-results/internal/utils"
)

// FailuresHandler holds the failures service.
type FailuresHandler struct {
	service service.FailuresServiceInterface
}

// NewFailuresHandler creates a new FailuresHandler.
func NewFailuresHandler(s service.FailuresServiceInterface) *FailuresHandler {
	return &FailuresHandler{service: s}
}

func (fh *FailuresHandler) GetBuildFailures(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.PathValue("id")
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID format: "+err.Error())
		return
	}

	failures, err := fh.service.GetFailuresByBuildID(buildID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching failures: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, failures)
}
