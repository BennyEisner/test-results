package handler

import (
	"database/sql"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/BennyEisner/test-results/internal/service"
	"github.com/BennyEisner/test-results/internal/utils"
)

// BuildInput defines the input for creating a build.
type BuildInput struct {
	BuildNumber string   `json:"build_number" xml:"build_number"`
	CIProvider  string   `json:"ci_provider" xml:"ci_provider"`
	CIURL       string   `json:"ci_url" xml:"ci_url"`
	Duration    *float64 `json:"duration" xml:"duration"`
}

// BuildCreateInput defines the input for creating a build with a test suite ID.
type BuildCreateInput struct {
	TestSuiteID int      `json:"test_suite_id" xml:"test_suite_id"`
	BuildNumber string   `json:"build_number" xml:"build_number"`
	CIProvider  string   `json:"ci_provider" xml:"ci_provider"`
	CIURL       string   `json:"ci_url" xml:"ci_url"`
	Duration    *float64 `json:"duration" xml:"duration"`
}

// BuildUpdateInput defines the input for updating a build.
type BuildUpdateInput struct {
	BuildNumber *string  `json:"build_number" xml:"build_number"`
	CIProvider  *string  `json:"ci_provider" xml:"ci_provider"`
	CIURL       *string  `json:"ci_url" xml:"ci_url"`
	Duration    *float64 `json:"duration" xml:"duration"`
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
	if m.Duration != nil {
		apiBuild.Duration = m.Duration
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

func (bh *BuildHandler) GetBuildsByTestSuiteID(w http.ResponseWriter, r *http.Request) {
	testSuiteIDStr := r.PathValue("suiteId")
	testSuiteID, err := strconv.ParseInt(testSuiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid test suite ID format: "+err.Error())
		return
	}

	builds, err := bh.service.GetBuildsByTestSuiteID(testSuiteID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching builds: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, toAPIBuilds(builds))
}

func (bh *BuildHandler) GetBuildByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID format: "+err.Error())
		return
	}

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

func (bh *BuildHandler) GetAllBuilds(w http.ResponseWriter, r *http.Request) {
	builds, err := bh.service.GetAllBuilds()
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching all builds: "+err.Error())
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, toAPIBuilds(builds))
}

func (bh *BuildHandler) GetRecentBuilds(w http.ResponseWriter, r *http.Request) {
	// Check if this is a project-specific request (from path parameter)
	projectIDStr := r.PathValue("id")
	if projectIDStr != "" {
		projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
		if err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format: "+err.Error())
			return
		}
		builds, err := bh.service.GetRecentBuildsByProjectID(projectID)
		if err != nil {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching recent builds for project: "+err.Error())
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, toAPIBuilds(builds))
	} else {
		// Check for query parameter as fallback
		projectIDStr = r.URL.Query().Get("projectId")
		if projectIDStr != "" {
			projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
			if err != nil {
				utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID format: "+err.Error())
				return
			}
			builds, err := bh.service.GetRecentBuildsByProjectID(projectID)
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching recent builds for project: "+err.Error())
				return
			}
			utils.RespondWithJSON(w, http.StatusOK, toAPIBuilds(builds))
		} else {
			builds, err := bh.service.GetAllBuilds()
			if err != nil {
				utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching recent builds: "+err.Error())
				return
			}
			utils.RespondWithJSON(w, http.StatusOK, toAPIBuilds(builds))
		}
	}
}

func (bh *BuildHandler) CreateBuildForTestSuite(w http.ResponseWriter, r *http.Request) {
	testSuiteIDStr := r.PathValue("suiteId")
	testSuiteID, err := strconv.ParseInt(testSuiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid test suite ID format: "+err.Error())
		return
	}

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
		Duration:    input.Duration,
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

// decodeBuildCreateInput decodes the request body into BuildCreateInput based on content type
func decodeBuildCreateInput(r *http.Request) (BuildCreateInput, error) {
	var input BuildCreateInput
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		if err := xml.NewDecoder(r.Body).Decode(&input); err != nil {
			return input, err
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			return input, err
		}
	}
	return input, nil
}

// validateBuildCreateInput validates the BuildCreateInput fields
func validateBuildCreateInput(input BuildCreateInput) error {
	if input.TestSuiteID == 0 {
		return errors.New("Test Suite ID is required and must be valid")
	}
	if strings.TrimSpace(input.BuildNumber) == "" {
		return errors.New("Build number is required")
	}
	if strings.TrimSpace(input.CIProvider) == "" {
		return errors.New("CI provider is required")
	}
	return nil
}

// checkTestSuiteExists checks if the test suite exists
func (bh *BuildHandler) checkTestSuiteExists(testSuiteID int64) error {
	exists, err := bh.service.CheckTestSuiteExists(testSuiteID)
	if err != nil {
		return fmt.Errorf("Database error checking test suite: %w", err)
	}
	if !exists {
		return fmt.Errorf("Test suite with ID %d not found", testSuiteID)
	}
	return nil
}

func (bh *BuildHandler) CreateBuild(w http.ResponseWriter, r *http.Request) {
	input, decodeErr := decodeBuildCreateInput(r)
	if decodeErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+decodeErr.Error())
		return
	}
	defer r.Body.Close()

	if err := validateBuildCreateInput(input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	testSuiteID64 := int64(input.TestSuiteID)
	if err := bh.checkTestSuiteExists(testSuiteID64); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.RespondWithError(w, http.StatusNotFound, err.Error())
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	modelBuild := models.Build{
		TestSuiteID: testSuiteID64,
		BuildNumber: input.BuildNumber,
		CIProvider:  input.CIProvider,
		Duration:    input.Duration,
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

func (bh *BuildHandler) DeleteBuild(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID format: "+err.Error())
		return
	}

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

// decodeBuildUpdateInput decodes the request body into BuildUpdateInput based on content type
func decodeBuildUpdateInput(r *http.Request) (BuildUpdateInput, error) {
	var input BuildUpdateInput
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
		if err := xml.NewDecoder(r.Body).Decode(&input); err != nil {
			return input, err
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			return input, err
		}
	}
	return input, nil
}

// validateBuildUpdateInput validates the BuildUpdateInput fields
func validateBuildUpdateInput(input BuildUpdateInput) error {
	if input.BuildNumber != nil && strings.TrimSpace(*input.BuildNumber) == "" {
		return errors.New("Build number cannot be empty if provided")
	}
	if input.CIProvider != nil && strings.TrimSpace(*input.CIProvider) == "" {
		return errors.New("CI provider cannot be empty if provided")
	}
	return nil
}

// respondUpdateBuildError handles error responses for UpdateBuild
func respondUpdateBuildError(w http.ResponseWriter, err error) {
	if err == sql.ErrNoRows {
		utils.RespondWithError(w, http.StatusNotFound, "Build not found")
	} else if strings.Contains(err.Error(), "no valid fields provided for update") {
		utils.RespondWithError(w, http.StatusBadRequest, "No valid fields provided for update")
	} else {
		utils.RespondWithError(w, http.StatusInternalServerError, "Update build failed: "+err.Error())
	}
}

func (bh *BuildHandler) UpdateBuild(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid build ID format: "+err.Error())
		return
	}
	defer r.Body.Close()

	input, decodeErr := decodeBuildUpdateInput(r)
	if decodeErr != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload: "+decodeErr.Error())
		return
	}

	if err := validateBuildUpdateInput(input); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedBuild, err := bh.service.UpdateBuild(id, input.BuildNumber, input.CIProvider, input.CIURL, input.Duration)
	if err != nil {
		respondUpdateBuildError(w, err)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, toAPIBuild(updatedBuild))
}

func (bh *BuildHandler) GetBuildDurationTrends(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.URL.Query().Get("projectId")
	if projectIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "projectId query parameter is required")
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid projectId format")
		return
	}

	suiteIDStr := r.URL.Query().Get("suiteId")
	if suiteIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "suiteId query parameter is required")
		return
	}

	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid suiteId format")
		return
	}

	trends, err := bh.service.GetBuildDurationTrends(projectID, suiteID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching build duration trends: "+err.Error())
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, trends)
}
