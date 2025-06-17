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

// HandleProjectTestSuites handles GET and POST requests for test suites associated with a specific project.
// Expected path: /api/projects/{projectID}/suites
func HandleProjectTestSuites(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Path segments for /api/projects/{projectID}/suites will be:
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/projects/"), "/")

	if len(pathSegments) < 1 || pathSegments[0] == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID in URL for suites")
		return
	}
	// Example: /api/projects/1/suites -> pathSegments = ["1", "suites"]
	// Example: /api/projects/1/suites/ -> pathSegments = ["1", "suites", ""]
	if len(pathSegments) < 2 || pathSegments[1] != "suites" {
		utils.RespondWithError(w, http.StatusBadRequest, "URL path for project suites is malformed, expected /api/projects/{projectID}/suites")
		return
	}

	projectIDStr := pathSegments[0]
	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format: "+err.Error())
		return
	}

	// Check if the project exists
	var projectExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1)", projectID).Scan(&projectExists)
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
		getTestSuitesByProjectID(w, r, projectID, db)
	case http.MethodPost:
		createTestSuiteForProject(w, r, projectID, db)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed for this endpoint")
	}
}

func getTestSuitesByProjectID(w http.ResponseWriter, r *http.Request, projectID int64, db *sql.DB) {
	rows, err := db.Query("SELECT id, project_id, name, parent_id, time FROM test_suites WHERE project_id = $1 ORDER BY name", projectID)
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
		if err := rows.Scan(&ts.ID, &ts.ProjectID, &ts.Name, &parentID, &ts.Time); err != nil {
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

func createTestSuiteForProject(w http.ResponseWriter, r *http.Request, projectID int64, db *sql.DB) {
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
		// Optional: Check if parent suite ID exists and belongs to the same project
		parentIDArg = sql.NullInt64{Int64: *input.ParentID, Valid: true}
	} else {
		parentIDArg = sql.NullInt64{Valid: false}
	}

	err := db.QueryRow(
		"INSERT INTO test_suites(project_id, name, parent_id, time) VALUES($1, $2, $3, $4) RETURNING id, project_id, name, parent_id, time",
		projectID, input.Name, parentIDArg, input.Time,
	).Scan(&createdSuite.ID, &createdSuite.ProjectID, &createdSuite.Name, &parentIDArg, &createdSuite.Time)

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

// GetProjectTestSuiteByID fetches a specific test suite by its ID and projectID.
func GetProjectTestSuiteByID(w http.ResponseWriter, r *http.Request, projectID int64, suiteID int64, db *sql.DB) {
	var ts models.TestSuite
	var parentID sql.NullInt64

	// First, verify the project exists to give a more specific error if project_id is wrong.
	var projectExists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1)", projectID).Scan(&projectExists)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error checking project existence: "+err.Error())
		return
	}
	if !projectExists {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Project with ID %d not found", projectID))
		return
	}

	// Now fetch the test suite, ensuring it belongs to the specified project.
	err = db.QueryRow("SELECT id, project_id, name, parent_id, time FROM test_suites WHERE id = $1 AND project_id = $2", suiteID, projectID).Scan(
		&ts.ID, &ts.ProjectID, &ts.Name, &parentID, &ts.Time)

	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Test suite with ID %d not found in project %d", suiteID, projectID))
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
	err := db.QueryRow("SELECT id, project_id, name, parent_id, time FROM test_suites WHERE id = $1", suiteID).Scan(
		&ts.ID, &ts.ProjectID, &ts.Name, &parentID, &ts.Time)
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
