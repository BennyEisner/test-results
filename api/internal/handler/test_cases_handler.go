package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/BennyEisner/test-results/internal/utils"
)

// TestCaseCreateInput defines the expected structure for creating a new test case definition.
type TestCaseCreateInput struct {
	Name      string `json:"name"`
	Classname string `json:"classname"`
	// Time and Status are not part of the definition anymore, they belong to BuildTestCaseExecution
}

// HandleSuiteTestCases handles GET and POST requests for test case definitions associated with a specific test suite.
// Expected path: /api/suites/{suiteId}/cases
func HandleSuiteTestCases(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/suites/"), "/")

	// Request error handling
	if len(pathSegments) < 1 || pathSegments[0] == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid suite ID in URL for test cases")
		return
	}
	if len(pathSegments) > 1 && pathSegments[1] != "cases" {
		utils.RespondWithError(w, http.StatusBadRequest, "URL path for suite test cases is malformed, expected /api/suites/{suiteId}/cases")
		return
	}
	suiteIDStr := pathSegments[0]
	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid suite ID format: "+err.Error())
		return
	}

	// Check if the suite exists
	var suiteExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM test_suites WHERE id = $1)", suiteID).Scan(&suiteExists)
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
		getTestCasesBySuiteID(w, r, suiteID, db)
	case http.MethodPost:
		createTestCaseForSuite(w, r, suiteID, db)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed for this endpoint")
	}
}

func getTestCasesBySuiteID(w http.ResponseWriter, r *http.Request, suiteID int64, db *sql.DB) {
	rows, err := db.Query("SELECT id, suite_id, name, classname FROM test_cases WHERE suite_id = $1 ORDER BY name", suiteID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching test case definitions: "+err.Error())
		return
	}
	defer rows.Close()

	cases := []models.TestCase{}
	for rows.Next() {
		var tc models.TestCase
		if err := rows.Scan(&tc.ID, &tc.SuiteID, &tc.Name, &tc.Classname); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning test case definition: "+err.Error())
			return
		}
		cases = append(cases, tc)
	}
	if err = rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating test case definition rows: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, cases)
}

func createTestCaseForSuite(w http.ResponseWriter, r *http.Request, suiteID int64, db *sql.DB) {
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

	// Status and Time are no longer part of creating a test case definition.
	// They are recorded with BuildTestCaseExecution.

	var createdCase models.TestCase
	err := db.QueryRow(
		"INSERT INTO test_cases(suite_id, name, classname) VALUES($1, $2, $3) RETURNING id, suite_id, name, classname",
		suiteID, input.Name, input.Classname,
	).Scan(&createdCase.ID, &createdCase.SuiteID, &createdCase.Name, &createdCase.Classname)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error creating test case definition: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, createdCase)
}

// HandleTestCaseByPath handles GET requests for a specific test case definition.
// Expected path: /api/cases/{caseId}
func HandleTestCaseByPath(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
		getTestCaseByID(w, r, caseID, db)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed for this endpoint")
	}
}

func getTestCaseByID(w http.ResponseWriter, r *http.Request, caseID int64, db *sql.DB) {
	var tc models.TestCase
	err := db.QueryRow("SELECT id, suite_id, name, classname FROM test_cases WHERE id = $1", caseID).Scan(
		&tc.ID, &tc.SuiteID, &tc.Name, &tc.Classname)
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
