package handler

import (
	"net/http"
	"strconv"
	"strings"

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

// HandleBuildFailures handles GET requests for failures associated with a specific build.
// Expected path: /api/builds/{buildId}/failures
func (fh *FailuresHandler) HandleBuildFailures(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only GET method is allowed")
		return
	}

	// Extract build ID from path
	// Expected path: /api/builds/{buildId}/failures
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/builds/"), "/")
	if len(pathSegments) != 2 || pathSegments[0] == "" || pathSegments[1] != "failures" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL path. Expected /api/builds/{buildId}/failures")
		return
	}

	buildIDStr := pathSegments[0]
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
