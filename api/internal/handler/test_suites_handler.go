package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BennyEisner/test-results/internal/service"
	"github.com/BennyEisner/test-results/internal/utils"
)

// TestSuiteCreateInput remains in the handler as it's specific to request decoding.
type TestSuiteCreateInput struct {
	Name     string  `json:"name"`
	ParentID *int64  `json:"parent_id,omitempty"`
	Time     float64 `json:"time"`
}

// TestSuiteHandler holds the test suite service.
type TestSuiteHandler struct {
	service service.TestSuiteServiceInterface
}

// NewTestSuiteHandler creates a new TestSuiteHandler.
func NewTestSuiteHandler(s service.TestSuiteServiceInterface) *TestSuiteHandler {
	return &TestSuiteHandler{service: s}
}

// HandleProjectTestSuites handles GET and POST requests for test suites associated with a specific project.
// Expected path: /api/projects/{projectID}/suites
func (tsh *TestSuiteHandler) HandleProjectTestSuites(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/projects/"), "/")

	if len(pathSegments) < 2 || pathSegments[0] == "" || pathSegments[1] != "suites" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL path for project suites. Expected /api/projects/{projectID}/suites")
		return
	}

	projectIDStr := pathSegments[0]
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format: "+err.Error())
		return
	}

	projectExists, err := tsh.service.CheckProjectExists(projectID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error checking project existence: "+err.Error())
		return
	}
	if !projectExists {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Project with ID %d not found", projectID))
		return
	}

	switch r.Method {
	case http.MethodGet:
		tsh.getTestSuitesByProjectID(w, r, projectID)
	case http.MethodPost:
		tsh.createTestSuiteForProject(w, r, projectID)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed for this endpoint")
	}
}

func (tsh *TestSuiteHandler) getTestSuitesByProjectID(w http.ResponseWriter, r *http.Request, projectID int64) {
	suites, err := tsh.service.GetTestSuitesByProjectID(projectID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching test suites: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, suites)
}

func (tsh *TestSuiteHandler) createTestSuiteForProject(w http.ResponseWriter, r *http.Request, projectID int64) {
	var input TestSuiteCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(input.Name) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Test suite name is required")
		return
	}

	// Optional: Validate parentID if provided
	if input.ParentID != nil {
		parentExists, err := tsh.service.CheckTestSuiteExists(*input.ParentID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error checking parent test suite %d: %s", *input.ParentID, err.Error()))
			return
		}
		if !parentExists {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Parent test suite with ID %d not found", *input.ParentID))
			return
		}
		// Further validation: ensure parent suite belongs to the same projectID.
		// This might require an additional service method or logic here.
		// For now, assuming DB constraints or later checks handle this.
	}

	createdSuite, err := tsh.service.CreateTestSuite(projectID, input.Name, input.ParentID, input.Time)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating test suite: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, createdSuite)
}

// GetProjectTestSuiteByID fetches a specific test suite by its ID and projectID.
// This function is called by the router, so it needs to be a method of TestSuiteHandler.
func (tsh *TestSuiteHandler) GetProjectTestSuiteByID(w http.ResponseWriter, r *http.Request, projectID int64, suiteID int64) {
	// Project existence check can be done first by the service or here.
	// The service method GetProjectTestSuiteByID implicitly handles this by querying with projectID.
	projectExists, err := tsh.service.CheckProjectExists(projectID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error checking project: "+err.Error())
		return
	}
	if !projectExists {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Project with ID %d not found", projectID))
		return
	}

	ts, err := tsh.service.GetProjectTestSuiteByID(projectID, suiteID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Test suite with ID %d not found in project %d", suiteID, projectID))
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching test suite: "+err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, ts)
}

// HandleTestSuiteByPath handles GET requests for a specific test suite by its ID only.
// Expected path: /api/suites/{suiteId}
func (tsh *TestSuiteHandler) HandleTestSuiteByPath(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/suites/"), "/")
	if len(pathSegments) < 1 || pathSegments[0] == "" || strings.Contains(pathSegments[0], "/") {
		// If pathSegments[0] contains "/", it means it's likely /api/suites/{id}/cases, which is handled elsewhere
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid test suite ID in URL. Expected /api/suites/{id}")
		return
	}

	suiteIDStr := pathSegments[0]
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid test suite ID format: "+err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		tsh.getTestSuiteByID(w, r, suiteID)
	// Add PUT, DELETE later if needed
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed for this endpoint")
	}
}

func (tsh *TestSuiteHandler) getTestSuiteByID(w http.ResponseWriter, r *http.Request, suiteID int64) {
	ts, err := tsh.service.GetTestSuiteByID(suiteID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Test suite not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching test suite: "+err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, ts)
}
