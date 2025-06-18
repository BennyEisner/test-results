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

// TestCaseCreateInput defines the expected structure for creating a new test case definition.
type TestCaseCreateInput struct {
	Name      string `json:"name"`
	Classname string `json:"classname"`
}

// TestCaseHandler holds the test case service.
type TestCaseHandler struct {
	service service.TestCaseServiceInterface
}

// NewTestCaseHandler creates a new TestCaseHandler.
func NewTestCaseHandler(s service.TestCaseServiceInterface) *TestCaseHandler {
	return &TestCaseHandler{service: s}
}

// HandleSuiteTestCases handles GET and POST requests for test case definitions associated with a specific test suite.
// Expected path: /api/suites/{suiteId}/cases
func (tch *TestCaseHandler) HandleSuiteTestCases(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/suites/"), "/")

	if len(pathSegments) < 2 || pathSegments[0] == "" || pathSegments[1] != "cases" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL path for suite test cases. Expected /api/suites/{suiteId}/cases")
		return
	}
	suiteIDStr := pathSegments[0]
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid suite ID format: "+err.Error())
		return
	}

	suiteExists, err := tch.service.CheckTestSuiteExists(suiteID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error checking suite existence: "+err.Error())
		return
	}
	if !suiteExists {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Test suite with ID %d not found", suiteID))
		return
	}

	switch r.Method {
	case http.MethodGet:
		tch.getTestCasesBySuiteID(w, r, suiteID)
	case http.MethodPost:
		tch.createTestCaseForSuite(w, r, suiteID)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed for this endpoint")
	}
}

func (tch *TestCaseHandler) getTestCasesBySuiteID(w http.ResponseWriter, r *http.Request, suiteID int64) {
	cases, err := tch.service.GetTestCasesBySuiteID(suiteID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching test case definitions: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, cases)
}

func (tch *TestCaseHandler) createTestCaseForSuite(w http.ResponseWriter, r *http.Request, suiteID int64) {
	var input TestCaseCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(input.Name) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Test case name is required")
		return
	}
	if strings.TrimSpace(input.Classname) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Test case classname is required")
		return
	}

	createdCase, err := tch.service.CreateTestCase(suiteID, input.Name, input.Classname)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error creating test case definition: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, createdCase)
}

// HandleTestCaseByPath handles GET requests for a specific test case definition.
// Expected path: /api/cases/{caseId}
func (tch *TestCaseHandler) HandleTestCaseByPath(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/cases/"), "/")
	if len(pathSegments) != 1 || pathSegments[0] == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid test case ID in URL")
		return
	}

	caseIDStr := pathSegments[0]
	caseID, err := strconv.ParseInt(caseIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid test case ID format: "+err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		tch.getTestCaseByID(w, r, caseID)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed for this endpoint")
	}
}

func (tch *TestCaseHandler) getTestCaseByID(w http.ResponseWriter, r *http.Request, caseID int64) {
	tc, err := tch.service.GetTestCaseByID(caseID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Test case definition not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching test case definition: "+err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, tc)
}
