package handler

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/BennyEisner/test-results/internal/service"
	"github.com/BennyEisner/test-results/internal/utils"
)

// BuildInput remains the same
type BuildInput struct {
	BuildNumber string `json:"build_number" xml:"build_number"`
	CIProvider  string `json:"ci_provider" xml:"ci_provider"`
	CIURL       string `json:"ci_url" xml:"ci_url"`
}

// BuildCreateInput remains the same
type BuildCreateInput struct {
	TestSuiteID int    `json:"test_suite_id" xml:"test_suite_id"`
	BuildNumber string `json:"build_number" xml:"build_number"`
	CIProvider  string `json:"ci_provider" xml:"ci_provider"`
	CIURL       string `json:"ci_url" xml:"ci_url"`
}

// BuildUpdateInput remains the same
type BuildUpdateInput struct {
	BuildNumber *string `json:"build_number" xml:"build_number"`
	CIProvider  *string `json:"ci_provider" xml:"ci_provider"`
	CIURL       *string `json:"ci_url" xml:"ci_url"`
}

// BuildHandler holds the build service.
type BuildHandler struct {
	service service.BuildServiceInterface
}

// NewBuildHandler creates a new BuildHandler.
func NewBuildHandler(s service.BuildServiceInterface) *BuildHandler {
	return &BuildHandler{service: s}
}

// Helper function to convert models.Build to utils.Build for API responses
func toAPIBuild(m *models.Build) utils.Build {
	apiBuild := utils.Build{
		ID:            int(m.ID),
		TestSuiteID:   int(m.TestSuiteID),
		ProjectID:     int(m.ProjectID),
		BuildNumber:   m.BuildNumber,
		CIProvider:    m.CIProvider,
		CreatedAt:     m.CreatedAt,
		TestCaseCount: int(m.TestCaseCount),
	}
	if m.CIURL != nil {
		apiBuild.CIURL = *m.CIURL
	}
	return apiBuild
}

// Helper function to convert a slice of models.Build to a slice of utils.Build
func toAPIBuilds(ms []models.Build) []utils.Build {
	apiBuilds := make([]utils.Build, len(ms))
	for i, m := range ms {
		apiBuilds[i] = toAPIBuild(&m) // Pass address of m as toAPIBuild expects a pointer
	}
	return apiBuilds
}

// HandleBuilds handles GET (all builds) and POST (create build) requests for /api/builds
func (bh *BuildHandler) HandleBuilds(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		bh.getAllBuilds(w, r)
	case http.MethodPost:
		bh.createBuild(w, r)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleBuildByPath handles operations on a specific build via /api/builds/{id}
func (bh *BuildHandler) HandleBuildByPath(w http.ResponseWriter, r *http.Request) {
	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/builds/"), "/")
	if len(pathSegments) != 1 || pathSegments[0] == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID in URL")
		return
	}

	idStr := pathSegments[0]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID format: "+err.Error())
		return
	}

	switch r.Method {
	case http.MethodGet:
		bh.getBuildByID(w, r, id)
	case http.MethodPatch:
		bh.updateBuild(w, r, id)
	case http.MethodDelete:
		bh.deleteBuild(w, r, id)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

// HandleTestSuiteBuilds handles GET and POST for builds under a specific test suite: /api/projects/{projectID}/test_suites/{testSuiteID}/builds
func (bh *BuildHandler) HandleTestSuiteBuilds(w http.ResponseWriter, r *http.Request) {
	trimmedPath := strings.TrimPrefix(r.URL.Path, "/api/projects/")
	pathSegments := strings.Split(strings.Trim(trimmedPath, "/"), "/")

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

	exists, err := bh.service.CheckTestSuiteExists(testSuiteID)
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
		bh.getBuildsByTestSuiteID(w, r, testSuiteID)
	case http.MethodPost:
		bh.createBuildForTestSuite(w, r, testSuiteID)
	default:
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
	}
}

func (bh *BuildHandler) getBuildsByTestSuiteID(w http.ResponseWriter, r *http.Request, testSuiteID int64) {
	builds, err := bh.service.GetBuildsByTestSuiteID(testSuiteID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching builds: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, toAPIBuilds(builds))
}

func (bh *BuildHandler) getBuildByID(w http.ResponseWriter, r *http.Request, id int64) {
	build, err := bh.service.GetBuildByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Build not found")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching build: "+err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, toAPIBuild(build))
}

func (bh *BuildHandler) getAllBuilds(w http.ResponseWriter, r *http.Request) {
	builds, err := bh.service.GetAllBuilds()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching all builds: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, toAPIBuilds(builds))
}

func (bh *BuildHandler) GetRecentBuilds(w http.ResponseWriter, r *http.Request) {
	builds, err := bh.service.GetAllBuilds()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching recent builds: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, toAPIBuilds(builds))
}

func (bh *BuildHandler) createBuildForTestSuite(w http.ResponseWriter, r *http.Request, testSuiteID int64) {
	var input BuildInput
	contentType := r.Header.Get("Content-Type")
	var decodeErr error

	if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		decodeErr = xml.NewDecoder(r.Body).Decode(&input)
	} else {
		decodeErr = json.NewDecoder(r.Body).Decode(&input)
	}
	defer r.Body.Close()

	if decodeErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+decodeErr.Error())
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

	modelBuild := models.Build{
		TestSuiteID: testSuiteID,
		BuildNumber: input.BuildNumber,
		CIProvider:  input.CIProvider,
	}
	if strings.TrimSpace(input.CIURL) != "" {
		modelBuild.CIURL = &input.CIURL
	}

	createdBuild, err := bh.service.CreateBuild(&modelBuild)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating build: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, toAPIBuild(createdBuild))
}

func (bh *BuildHandler) createBuild(w http.ResponseWriter, r *http.Request) {
	var input BuildCreateInput
	contentType := r.Header.Get("Content-Type")
	var decodeErr error

	if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		decodeErr = xml.NewDecoder(r.Body).Decode(&input)
	} else {
		decodeErr = json.NewDecoder(r.Body).Decode(&input)
	}
	defer r.Body.Close()

	if decodeErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+decodeErr.Error())
		return
	}

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

	testSuiteID64 := int64(input.TestSuiteID)
	exists, err := bh.service.CheckTestSuiteExists(testSuiteID64)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Database error checking test suite: "+err.Error())
		return
	}
	if !exists {
		utils.RespondWithError(w, http.StatusNotFound, fmt.Sprintf("Test suite with ID %d not found", input.TestSuiteID))
		return
	}

	modelBuild := models.Build{
		TestSuiteID: testSuiteID64,
		BuildNumber: input.BuildNumber,
		CIProvider:  input.CIProvider,
	}
	if strings.TrimSpace(input.CIURL) != "" {
		modelBuild.CIURL = &input.CIURL
	}

	createdBuild, err := bh.service.CreateBuild(&modelBuild)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating build: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, toAPIBuild(createdBuild))
}

func (bh *BuildHandler) deleteBuild(w http.ResponseWriter, r *http.Request, id int64) {
	rowsAffected, err := bh.service.DeleteBuild(id)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting build: "+err.Error())
		return
	}
	if rowsAffected == 0 {
		utils.RespondWithError(w, http.StatusNotFound, "Build not found or already deleted")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Build deleted successfully"})
}

func (bh *BuildHandler) updateBuild(w http.ResponseWriter, r *http.Request, id int64) {
	var input BuildUpdateInput
	contentType := r.Header.Get("Content-Type")
	var decodeErr error

	if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		decodeErr = xml.NewDecoder(r.Body).Decode(&input)
	} else {
		decodeErr = json.NewDecoder(r.Body).Decode(&input)
	}
	defer r.Body.Close()

	if decodeErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+decodeErr.Error())
		return
	}

	// Input validation for non-nil fields
	if input.BuildNumber != nil && strings.TrimSpace(*input.BuildNumber) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "Build number cannot be empty if provided")
		return
	}
	if input.CIProvider != nil && strings.TrimSpace(*input.CIProvider) == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "CI provider cannot be empty if provided")
		return
	}
	// CIURL can be an empty string to clear it, so no trim check here for emptiness,
	// but the service layer handles sql.NullString for empty vs. nil.

	updatedBuild, err := bh.service.UpdateBuild(id, input.BuildNumber, input.CIProvider, input.CIURL)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.RespondWithError(w, http.StatusNotFound, "Build not found")
		} else if strings.Contains(err.Error(), "no valid fields provided for update") { // Check for specific service error
			utils.RespondWithError(w, http.StatusBadRequest, "No valid fields provided for update")
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Update build failed: "+err.Error())
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, toAPIBuild(updatedBuild))
}
