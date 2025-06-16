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

// TestSuiteCreateInput is the expected structure for creating a new test suite.
type TestSuiteCreateInput struct {
	Name     string  `json:"name"`
	ParentID *int64  `json:"parent_id,omitempty"`
	Time     float64 `json:"time"`
}

// HandleBuildTestSuites handles GET and POST requests for test suites associated with a specific build.
// Expected path: /api/builds/{buildId}/suites
func HandleBuildTestSuites(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Path segments for /api/builds/{buildId}/suites will be:
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/builds/"), "/")

	if len(pathSegments) < 1 || pathSegments[0] == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID in URL for suites")
		return
	}
	if len(pathSegments) > 1 && pathSegments[1] != "suites" {
		// This case should ideally be caught by router logic, but good to double check.
		utils.RespondWithError(w, http.StatusBadRequest, "URL path for build suites is malformed, expected /api/builds/{buildId}/suites")
		return
	}

	buildIDStr := pathSegments[0]
	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID format: "+err.Error())
		return
	}

	// Check if the build exists
	var buildExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM builds WHERE id = $1)", buildID).Scan(&buildExists)
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
		getTestSuitesByBuildID(w, r, buildID, db)
	case http.MethodPost:
		createTestSuiteForBuild(w, r, buildID, db)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed for this endpoint")
	}
}

func getTestSuitesByBuildID(w http.ResponseWriter, r *http.Request, buildID int64, db *sql.DB) {
	rows, err := db.Query("SELECT id, build_id, name, parent_id, time FROM test_suites WHERE build_id = $1 ORDER BY name", buildID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching test suites: "+err.Error())
		return
	}
	defer rows.Close()

	suites := []models.TestSuite{}
	for rows.Next() {
		var ts models.TestSuite
		// Need to handle nullable ParentID
		var parentID sql.NullInt64
		if err := rows.Scan(&ts.ID, &ts.BuildID, &ts.Name, &parentID, &ts.Time); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning test suite: "+err.Error())
			return
		}
		if parentID.Valid {
			ts.ParentID = &parentID.Int64
		} else {
			ts.ParentID = nil
		}
		suites = append(suites, ts)
	}
	if err = rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating test suite rows: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, suites)
}

func createTestSuiteForBuild(w http.ResponseWriter, r *http.Request, buildID int64, db *sql.DB) {
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
	// Time could be 0, so not checking for > 0 strictly unless it's a business rule.

	var createdSuite models.TestSuite
	var parentIDArg sql.NullInt64
	if input.ParentID != nil {
		// Optional: Check if parent suite ID exists and belongs to the same build
		parentIDArg = sql.NullInt64{Int64: *input.ParentID, Valid: true}
	} else {
		parentIDArg = sql.NullInt64{Valid: false}
	}

	err := db.QueryRow(
		"INSERT INTO test_suites(build_id, name, parent_id, time) VALUES($1, $2, $3, $4) RETURNING id, build_id, name, parent_id, time",
		buildID, input.Name, parentIDArg, input.Time,
	).Scan(&createdSuite.ID, &createdSuite.BuildID, &createdSuite.Name, &parentIDArg, &createdSuite.Time)

	if err != nil {
		// More specific error for foreign key violation on parent_id if possible
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error creating test suite: "+err.Error())
		return
	}
	if parentIDArg.Valid {
		createdSuite.ParentID = &parentIDArg.Int64
	} else {
		createdSuite.ParentID = nil
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdSuite)
}

// HandleTestSuiteByPath handles GET requests for a specific test suite.
// Expected path: /api/suites/{suiteId}
func HandleTestSuiteByPath(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/suites/"), "/")
	if len(pathSegments) != 1 || pathSegments[0] == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid test suite ID in URL")
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
		getTestSuiteByID(w, r, suiteID, db)
	// Add PUT, DELETE later if needed
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed for this endpoint")
	}
}

func getTestSuiteByID(w http.ResponseWriter, r *http.Request, suiteID int64, db *sql.DB) {
	var ts models.TestSuite
	var parentID sql.NullInt64
	err := db.QueryRow("SELECT id, build_id, name, parent_id, time FROM test_suites WHERE id = $1", suiteID).Scan(
		&ts.ID, &ts.BuildID, &ts.Name, &parentID, &ts.Time)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Test suite not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching test suite: "+err.Error())
		}
		return
	}
	if parentID.Valid {
		ts.ParentID = &parentID.Int64
	} else {
		ts.ParentID = nil
	}
	utils.RespondWithJSON(w, http.StatusOK, ts)
}
