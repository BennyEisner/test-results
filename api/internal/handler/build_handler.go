package handler

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/BennyEisner/test-results/internal/models" // Using models for DB interaction
	"github.com/BennyEisner/test-results/internal/utils"
)

// First, define XML-compatible structs that will be used for input parsing
type BuildInput struct {
	BuildNumber string `json:"build_number" xml:"build_number"`
	CIProvider  string `json:"ci_provider" xml:"ci_provider"`
	CIURL       string `json:"ci_url" xml:"ci_url"`
}

type BuildCreateInput struct {
	TestSuiteID int    `json:"test_suite_id" xml:"test_suite_id"`
	BuildNumber string `json:"build_number" xml:"build_number"`
	CIProvider  string `json:"ci_provider" xml:"ci_provider"`
	CIURL       string `json:"ci_url" xml:"ci_url"`
}

type BuildUpdateInput struct {
	BuildNumber *string `json:"build_number" xml:"build_number"`
	CIProvider  *string `json:"ci_provider" xml:"ci_provider"`
	CIURL       *string `json:"ci_url" xml:"ci_url"`
}

// HandleBuilds handles GET (all builds) and POST (create build) requests for /api/builds
func HandleBuilds(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	switch r.Method {
	case http.MethodGet:
		getAllBuilds(w, r, db) // Renamed from GetBuilds
	case http.MethodPost:
		createBuild(w, r, db) // New function to handle general build creation
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleBuildByPath handles operations on a specific build via /api/builds/{id}
func HandleBuildByPath(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/builds/"), "/")
	if len(pathSegments) != 1 || pathSegments[0] == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID in URL")
		return
	}

	idStr := pathSegments[0]
	id, err := strconv.ParseInt(idStr, 10, 64) // Build ID is int64 in models.Build
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID format: "+err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		getBuildByID(w, r, id, db)
	case http.MethodPatch:
		updateBuild(w, r, id, db)
	case http.MethodDelete:
		deleteBuild(w, r, id, db)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleTestSuiteBuilds handles GET and POST for builds under a specific test suite: /api/projects/{projectID}/test_suites/{testSuiteID}/builds
func HandleTestSuiteBuilds(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Example path: /api/projects/1/suites/2/builds
	trimmedPath := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	// pathSegments for "1/suites/2/builds" will be ["1", "suites", "2", "builds"]
	pathSegments := strings.Split(strings.Trim(trimmedPath, "/"), "/")

	// Expected segments: {projectID}, "suites", {testSuiteID}, "builds"
	if len(pathSegments) < 4 || pathSegments[1] != "suites" || pathSegments[3] != "builds" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL. Expected /api/projects/{projectID}/suites/{testSuiteID}/builds")
		return
	}

	testSuiteIDStr := pathSegments[2]
	testSuiteID, err := strconv.ParseInt(testSuiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid test suite ID format: "+err.Error())
		return
	}

	// Check if the test suite exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM test_suites WHERE id = $1)", testSuiteID).Scan(&exists)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error checking test suite: "+err.Error())
		return
	}
	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Test suite with ID %d not found", testSuiteID))
		return
	}

	switch r.Method {
	case http.MethodGet:
		getBuildsByTestSuiteID(w, r, testSuiteID, db)
	case http.MethodPost:
		createBuildForTestSuite(w, r, testSuiteID, db)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// getBuildsByTestSuiteID fetches all builds for a given testSuiteID
func getBuildsByTestSuiteID(w http.ResponseWriter, r *http.Request, testSuiteID int64, db *sql.DB) {
	rows, err := db.Query("SELECT id, test_suite_id, build_number, ci_provider, ci_url, created_at FROM builds WHERE test_suite_id = $1 ORDER BY created_at DESC", testSuiteID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching builds: "+err.Error())
		return
	}
	defer rows.Close()

	buildsAPI := []utils.Build{} // Slice of utils.Build for API response
	for rows.Next() {
		var b models.Build // Scan into models.Build to handle db types
		var ciURL sql.NullString
		if err := rows.Scan(&b.ID, &b.TestSuiteID, &b.BuildNumber, &b.CIProvider, &ciURL, &b.CreatedAt); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning build: "+err.Error())
			return
		}
		apiBuild := utils.Build{
			ID:          int(b.ID),          // Convert int64 to int for API DTO
			TestSuiteID: int(b.TestSuiteID), // Convert int64 to int
			BuildNumber: b.BuildNumber,
			CIProvider:  b.CIProvider,
			CreatedAt:   b.CreatedAt,
		}
		if ciURL.Valid {
			apiBuild.CIURL = ciURL.String
		}
		buildsAPI = append(buildsAPI, apiBuild)
	}
	if err = rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating build rows: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, buildsAPI)
}

// getBuildByID fetches a single build by its ID
func getBuildByID(w http.ResponseWriter, r *http.Request, id int64, db *sql.DB) {
	var b models.Build
	var ciURL sql.NullString
	err := db.QueryRow("SELECT id, test_suite_id, build_number, ci_provider, ci_url, created_at FROM builds WHERE id = $1", id).Scan(
		&b.ID, &b.TestSuiteID, &b.BuildNumber, &b.CIProvider, &ciURL, &b.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Build not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching build: "+err.Error())
		}
		return
	}
	apiBuild := utils.Build{
		ID:          int(b.ID),
		TestSuiteID: int(b.TestSuiteID),
		BuildNumber: b.BuildNumber,
		CIProvider:  b.CIProvider,
		CreatedAt:   b.CreatedAt,
	}
	if ciURL.Valid {
		apiBuild.CIURL = ciURL.String
	}
	utils.RespondWithJSON(w, http.StatusOK, apiBuild)
}

// getAllBuilds fetches all builds from the database
func getAllBuilds(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	rows, err := db.Query("SELECT id, test_suite_id, build_number, ci_provider, ci_url, created_at FROM builds ORDER BY created_at DESC")
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error fetching all builds: "+err.Error())
		return
	}
	defer rows.Close()

	buildsAPI := []utils.Build{}
	for rows.Next() {
		var b models.Build
		var ciURL sql.NullString
		if err := rows.Scan(&b.ID, &b.TestSuiteID, &b.BuildNumber, &b.CIProvider, &ciURL, &b.CreatedAt); err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error scanning build: "+err.Error())
			return
		}
		apiBuild := utils.Build{
			ID:          int(b.ID),
			TestSuiteID: int(b.TestSuiteID),
			BuildNumber: b.BuildNumber,
			CIProvider:  b.CIProvider,
			CreatedAt:   b.CreatedAt,
		}
		if ciURL.Valid {
			apiBuild.CIURL = ciURL.String
		}
		buildsAPI = append(buildsAPI, apiBuild)
	}
	if err = rows.Err(); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error iterating build rows: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, buildsAPI)
}

// Modified to support XML input
func createBuildForTestSuite(w http.ResponseWriter, r *http.Request, testSuiteID int64, db *sql.DB) {
	var input BuildInput

	contentType := r.Header.Get("Content-Type")
	var decodeErr error

	if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		decodeErr = xml.NewDecoder(r.Body).Decode(&input)
	} else {
		// Default to JSON for other content types or if not specified
		decodeErr = json.NewDecoder(r.Body).Decode(&input)
	}

	if decodeErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+decodeErr.Error())
		return
	}
	defer r.Body.Close()

	if strings.TrimSpace(input.BuildNumber) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Build number is required")
		return
	}
	if strings.TrimSpace(input.CIProvider) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "CI provider is required")
		return
	}

	var newBuildID int64
	var createdAt time.Time
	ciURLNullStr := sql.NullString{String: input.CIURL, Valid: strings.TrimSpace(input.CIURL) != ""}

	err := db.QueryRow(
		"INSERT INTO builds(test_suite_id, build_number, ci_provider, ci_url, created_at) VALUES($1, $2, $3, $4, NOW()) RETURNING id, created_at",
		testSuiteID, input.BuildNumber, input.CIProvider, ciURLNullStr,
	).Scan(&newBuildID, &createdAt)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error creating build: "+err.Error())
		return
	}

	createdAPIBuild := utils.Build{
		ID:          int(newBuildID),
		TestSuiteID: int(testSuiteID), // testSuiteID from path
		BuildNumber: input.BuildNumber,
		CIProvider:  input.CIProvider,
		CIURL:       input.CIURL, // Return the string version
		CreatedAt:   createdAt,
	}
	utils.RespondWithJSON(w, http.StatusCreated, createdAPIBuild)
}

// Modified to support XML input
func createBuild(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var input BuildCreateInput

	contentType := r.Header.Get("Content-Type")
	var decodeErr error

	if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		decodeErr = xml.NewDecoder(r.Body).Decode(&input)
	} else {
		// Default to JSON
		decodeErr = json.NewDecoder(r.Body).Decode(&input)
	}

	if decodeErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+decodeErr.Error())
		return
	}
	defer r.Body.Close()

	if input.TestSuiteID == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Test Suite ID is required and must be valid")
		return
	}
	if strings.TrimSpace(input.BuildNumber) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Build number is required")
		return
	}
	if strings.TrimSpace(input.CIProvider) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "CI provider is required")
		return
	}

	// Check if the referenced test suite exists
	var testSuiteExists bool
	testSuiteID64 := int64(input.TestSuiteID) // Convert API int to model int64
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM test_suites WHERE id = $1)", testSuiteID64).Scan(&testSuiteExists)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error checking test suite: "+err.Error())
		return
	}
	if !testSuiteExists {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Test suite with ID %d not found", input.TestSuiteID))
		return
	}

	var newBuildID int64
	var createdAt time.Time
	ciURLNullStr := sql.NullString{String: input.CIURL, Valid: strings.TrimSpace(input.CIURL) != ""}

	err = db.QueryRow(
		"INSERT INTO builds(test_suite_id, build_number, ci_provider, ci_url, created_at) VALUES($1, $2, $3, $4, NOW()) RETURNING id, created_at",
		testSuiteID64, input.BuildNumber, input.CIProvider, ciURLNullStr,
	).Scan(&newBuildID, &createdAt)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error creating build: "+err.Error())
		return
	}

	createdAPIBuild := utils.Build{
		ID:          int(newBuildID),
		TestSuiteID: input.TestSuiteID,
		BuildNumber: input.BuildNumber,
		CIProvider:  input.CIProvider,
		CIURL:       input.CIURL,
		CreatedAt:   createdAt,
	}
	utils.RespondWithJSON(w, http.StatusCreated, createdAPIBuild)
}

// deleteBuild deletes a build by its ID
func deleteBuild(w http.ResponseWriter, r *http.Request, id int64, db *sql.DB) {
	result, err := db.Exec("DELETE FROM builds WHERE id = $1", id) // Target builds table
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error deleting build: "+err.Error())
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error checking delete result: "+err.Error())
		return
	}

	if rowsAffected == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Build not found or already deleted")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Build deleted successfully"})
}

// Supports xml and json
func updateBuild(w http.ResponseWriter, r *http.Request, id int64, db *sql.DB) {
	var input BuildUpdateInput

	contentType := r.Header.Get("Content-Type")
	var decodeErr error

	if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		decodeErr = xml.NewDecoder(r.Body).Decode(&input)
	} else {
		// Default to JSON
		decodeErr = json.NewDecoder(r.Body).Decode(&input)
	}

	if decodeErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+decodeErr.Error())
		return
	}
	defer r.Body.Close()

	// Check if build exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM builds WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error checking build: "+err.Error())
		return
	}
	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, "Build not found")
		return
	}

	updateFields := []string{}
	args := []interface{}{}
	argID := 1

	if input.BuildNumber != nil {
		if strings.TrimSpace(*input.BuildNumber) == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Build number cannot be empty if provided")
			return
		}
		updateFields = append(updateFields, fmt.Sprintf("build_number = $%d", argID))
		args = append(args, *input.BuildNumber)
		argID++
	}
	if input.CIProvider != nil {
		if strings.TrimSpace(*input.CIProvider) == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "CI provider cannot be empty if provided")
			return
		}
		updateFields = append(updateFields, fmt.Sprintf("ci_provider = $%d", argID))
		args = append(args, *input.CIProvider)
		argID++
	}
	if input.CIURL != nil { // If CIURL key is present in JSON/XML
		ciURLToUpdate := sql.NullString{String: *input.CIURL, Valid: strings.TrimSpace(*input.CIURL) != ""}
		updateFields = append(updateFields, fmt.Sprintf("ci_url = $%d", argID))
		args = append(args, ciURLToUpdate)
		argID++
	}

	if len(updateFields) == 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "No valid fields provided for update")
		return
	}

	args = append(args, id) // Add ID for WHERE clause
	query := fmt.Sprintf("UPDATE builds SET %s WHERE id = $%d RETURNING id, test_suite_id, build_number, ci_provider, ci_url, created_at",
		strings.Join(updateFields, ", "), argID)

	var updatedBuildModel models.Build
	var updatedCIURL sql.NullString
	err = db.QueryRow(query, args...).Scan(
		&updatedBuildModel.ID,
		&updatedBuildModel.TestSuiteID,
		&updatedBuildModel.BuildNumber,
		&updatedBuildModel.CIProvider,
		&updatedCIURL,
		&updatedBuildModel.CreatedAt,
	)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Update build failed: "+err.Error())
		return
	}

	responseAPIBuild := utils.Build{
		ID:          int(updatedBuildModel.ID),
		TestSuiteID: int(updatedBuildModel.TestSuiteID),
		BuildNumber: updatedBuildModel.BuildNumber,
		CIProvider:  updatedBuildModel.CIProvider,
		CreatedAt:   updatedBuildModel.CreatedAt,
	}
	if updatedCIURL.Valid {
		responseAPIBuild.CIURL = updatedCIURL.String
	}

	utils.RespondWithJSON(w, http.StatusOK, responseAPIBuild)
}
