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

func (tsh *TestSuiteHandler) GetTestSuitesByProjectID(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format: "+err.Error())
		return
	}
	name := r.URL.Query().Get("name")
	if name != "" {
		suite, err := tsh.service.GetTestSuiteByName(projectID, name)
		if err != nil {
			if err == sql.ErrNoRows {
				utils.RespondWithError(w, http.StatusNotFound, "Suite not found")
			} else {
				utils.RespondWithError(w, http.StatusInternalServerError, "Database error: "+err.Error())
			}
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, suite)
		return
	}
	suites, err := tsh.service.GetTestSuitesByProjectID(projectID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching test suites: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, suites)
}

func (tsh *TestSuiteHandler) CreateTestSuiteForProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("id")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format: "+err.Error())
		return
	}

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
	}

	createdSuite, err := tsh.service.CreateTestSuite(projectID, input.Name, input.ParentID, input.Time)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating test suite: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, createdSuite)
}

func (tsh *TestSuiteHandler) GetProjectTestSuiteByID(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("projectId")
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format: "+err.Error())
		return
	}

	suiteIDStr := r.PathValue("suiteId")
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid suite ID format: "+err.Error())
		return
	}

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

func (tsh *TestSuiteHandler) GetTestSuiteByID(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.PathValue("id")
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid test suite ID format: "+err.Error())
		return
	}

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
