package handler

import (
	"database/sql"
	"encoding/json"
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

func (tch *TestCaseHandler) GetSuiteTestCases(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.PathValue("id")
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid suite ID format: "+err.Error())
		return
	}

	cases, err := tch.service.GetTestCasesBySuiteID(suiteID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching test case definitions: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, cases)
}

func (tch *TestCaseHandler) CreateTestCaseForSuite(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.PathValue("id")
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid suite ID format: "+err.Error())
		return
	}

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

func (tch *TestCaseHandler) GetTestCaseByID(w http.ResponseWriter, r *http.Request) {
	caseIDStr := r.PathValue("id")
	caseID, err := strconv.ParseInt(caseIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid test case ID format: "+err.Error())
		return
	}

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
