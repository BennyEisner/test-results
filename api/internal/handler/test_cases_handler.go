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

// TestCaseCreateInput defines the expected structure for creating a new test case.
type TestCaseCreateInput struct {
	Name      string  `json:"name"`
	Classname string  `json:"classname"`
	Time      float64 `json:"time"`
	Status    string  `json:"status"` // "passed", "failed", "skipped"
}

// HandleSuiteTestCases handles GET and POST requests for test cases associated with a specific test suite.
// Expected path: /api/suites/{suiteId}/cases
func HandleSuiteTestCases(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/suites/"), "/")

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
	rows, err := db.Query("SELECT id, suite_id, name, classname, time, status FROM test_cases WHERE suite_id = $1 ORDER BY name", suiteID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching test cases: "+err.Error())
		return
	}
	defer rows.Close()

	cases := []models.TestCase{}
	for rows.Next() {
		var tc models.TestCase
		if err := rows.Scan(&tc.ID, &tc.SuiteID, &tc.Name, &tc.Classname, &tc.Time, &tc.Status); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning test case: "+err.Error())
			return
		}
		cases = append(cases, tc)
	}
	if err = rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating test case rows: "+err.Error())
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

	// Validate status if provided
	validStatuses := map[string]bool{"passed": true, "failed": true, "skipped": true}
	statusToInsert := strings.ToLower(input.Status)
	if statusToInsert == "" {
		statusToInsert = "passed" // Use default if empty
	} else if !validStatuses[statusToInsert] {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid status. Must be one of: passed, failed, skipped")
		return
	}

	var createdCase models.TestCase
	err := db.QueryRow(
		"INSERT INTO test_cases(suite_id, name, classname, time, status) VALUES($1, $2, $3, $4, $5) RETURNING id, suite_id, name, classname, time, status",
		suiteID, input.Name, input.Classname, input.Time, statusToInsert,
	).Scan(&createdCase.ID, &createdCase.SuiteID, &createdCase.Name, &createdCase.Classname, &createdCase.Time, &createdCase.Status)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error creating test case: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, createdCase)
}

// HandleTestCaseByPath handles GET requests for a specific test case.
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
	err := db.QueryRow("SELECT id, suite_id, name, classname, time, status FROM test_cases WHERE id = $1", caseID).Scan(
		&tc.ID, &tc.SuiteID, &tc.Name, &tc.Classname, &tc.Time, &tc.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Test case not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching test case: "+err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, tc)
}
