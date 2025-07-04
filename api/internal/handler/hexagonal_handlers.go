package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/BennyEisner/test-results/internal/domain"
	"github.com/BennyEisner/test-results/internal/utils"
)

// HexagonalBuildHandler handles HTTP requests for build operations using hexagonal architecture
type HexagonalBuildHandler struct {
	buildService domain.BuildService
}

func NewHexagonalBuildHandler(buildService domain.BuildService) *HexagonalBuildHandler {
	return &HexagonalBuildHandler{buildService: buildService}
}

func (h *HexagonalBuildHandler) GetAllBuilds(w http.ResponseWriter, r *http.Request) {
	// This would need to be implemented based on business requirements
	// For now, return empty array as placeholder
	utils.RespondWithJSON(w, http.StatusOK, []interface{}{})
}

func (h *HexagonalBuildHandler) CreateBuild(w http.ResponseWriter, r *http.Request) {
	var build domain.Build
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	createdBuild, err := h.buildService.CreateBuild(r.Context(), &build)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to create build")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, createdBuild)
}

func (h *HexagonalBuildHandler) GetBuildByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing build ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid build ID format")
		return
	}

	build, err := h.buildService.GetBuildByID(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrBuildNotFound:
			utils.RespondWithError(w, http.StatusNotFound, "build not found")
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, build)
}

func (h *HexagonalBuildHandler) UpdateBuild(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing build ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid build ID format")
		return
	}

	var build domain.Build
	if err := json.NewDecoder(r.Body).Decode(&build); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	updatedBuild, err := h.buildService.UpdateBuild(r.Context(), id, &build)
	if err != nil {
		switch err {
		case domain.ErrBuildNotFound:
			utils.RespondWithError(w, http.StatusNotFound, "build not found")
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to update build")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, updatedBuild)
}

func (h *HexagonalBuildHandler) DeleteBuild(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing build ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid build ID format")
		return
	}

	if err := h.buildService.DeleteBuild(r.Context(), id); err != nil {
		switch err {
		case domain.ErrBuildNotFound:
			utils.RespondWithError(w, http.StatusNotFound, "build not found")
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to delete build")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "build deleted successfully"})
}

func (h *HexagonalBuildHandler) GetBuildsByProjectID(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("id")
	if projectIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing project ID")
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid project ID format")
		return
	}

	builds, err := h.buildService.GetBuildsByProjectID(r.Context(), projectID)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to get builds")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, builds)
}

func (h *HexagonalBuildHandler) GetBuildsByTestSuiteID(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.PathValue("suiteId")
	if suiteIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing test suite ID")
		return
	}

	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid test suite ID format")
		return
	}

	builds, err := h.buildService.GetBuildsByTestSuiteID(r.Context(), suiteID)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to get builds")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, builds)
}

// HexagonalTestSuiteHandler handles HTTP requests for test suite operations using hexagonal architecture
type HexagonalTestSuiteHandler struct {
	testSuiteService domain.TestSuiteService
}

func NewHexagonalTestSuiteHandler(testSuiteService domain.TestSuiteService) *HexagonalTestSuiteHandler {
	return &HexagonalTestSuiteHandler{testSuiteService: testSuiteService}
}

func (h *HexagonalTestSuiteHandler) GetTestSuiteByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing test suite ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid test suite ID format")
		return
	}

	testSuite, err := h.testSuiteService.GetTestSuiteByID(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrTestSuiteNotFound:
			utils.RespondWithError(w, http.StatusNotFound, "test suite not found")
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, testSuite)
}

func (h *HexagonalTestSuiteHandler) GetTestSuitesByProjectID(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("id")
	if projectIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing project ID")
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid project ID format")
		return
	}

	testSuites, err := h.testSuiteService.GetTestSuitesByProjectID(r.Context(), projectID)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to get test suites")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, testSuites)
}

func (h *HexagonalTestSuiteHandler) CreateTestSuiteForProject(w http.ResponseWriter, r *http.Request) {
	projectIDStr := r.PathValue("id")
	if projectIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing project ID")
		return
	}

	projectID, err := strconv.ParseInt(projectIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid project ID format")
		return
	}

	var request struct {
		Name     string `json:"name"`
		ParentID *int64 `json:"parent_id,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	testSuite, err := h.testSuiteService.CreateTestSuite(r.Context(), projectID, request.Name, request.ParentID)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		case domain.ErrDuplicateTestSuite:
			utils.RespondWithError(w, http.StatusConflict, "test suite with this name already exists")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to create test suite")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, testSuite)
}

func (h *HexagonalTestSuiteHandler) GetProjectTestSuiteByID(w http.ResponseWriter, r *http.Request) {
	// This would need to be implemented based on business requirements
	// For now, delegate to GetTestSuiteByID
	h.GetTestSuiteByID(w, r)
}

// HexagonalTestCaseHandler handles HTTP requests for test case operations using hexagonal architecture
type HexagonalTestCaseHandler struct {
	testCaseService domain.TestCaseService
}

func NewHexagonalTestCaseHandler(testCaseService domain.TestCaseService) *HexagonalTestCaseHandler {
	return &HexagonalTestCaseHandler{testCaseService: testCaseService}
}

func (h *HexagonalTestCaseHandler) GetTestCaseByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing test case ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid test case ID format")
		return
	}

	testCase, err := h.testCaseService.GetTestCaseByID(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrTestCaseNotFound:
			utils.RespondWithError(w, http.StatusNotFound, "test case not found")
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, testCase)
}

func (h *HexagonalTestCaseHandler) GetSuiteTestCases(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.PathValue("id")
	if suiteIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing suite ID")
		return
	}

	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid suite ID format")
		return
	}

	testCases, err := h.testCaseService.GetTestCasesBySuiteID(r.Context(), suiteID)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to get test cases")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, testCases)
}

func (h *HexagonalTestCaseHandler) CreateTestCaseForSuite(w http.ResponseWriter, r *http.Request) {
	suiteIDStr := r.PathValue("id")
	if suiteIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing suite ID")
		return
	}

	suiteID, err := strconv.ParseInt(suiteIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid suite ID format")
		return
	}

	var request struct {
		Name      string `json:"name"`
		Classname string `json:"classname"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	testCase, err := h.testCaseService.CreateTestCase(r.Context(), suiteID, request.Name, request.Classname)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		case domain.ErrDuplicateTestCase:
			utils.RespondWithError(w, http.StatusConflict, "test case with this name already exists")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to create test case")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, testCase)
}

func (h *HexagonalTestCaseHandler) GetMostFailedTests(w http.ResponseWriter, r *http.Request) {
	// This would need to be implemented based on business requirements
	// For now, return empty array as placeholder
	utils.RespondWithJSON(w, http.StatusOK, []interface{}{})
}

// HexagonalBuildExecutionHandler handles HTTP requests for build execution operations using hexagonal architecture
type HexagonalBuildExecutionHandler struct {
	buildExecutionService domain.BuildTestCaseExecutionService
}

func NewHexagonalBuildExecutionHandler(buildExecutionService domain.BuildTestCaseExecutionService) *HexagonalBuildExecutionHandler {
	return &HexagonalBuildExecutionHandler{buildExecutionService: buildExecutionService}
}

func (h *HexagonalBuildExecutionHandler) GetBuildExecutions(w http.ResponseWriter, r *http.Request) {
	buildIDStr := r.PathValue("id")
	if buildIDStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing build ID")
		return
	}

	buildID, err := strconv.ParseInt(buildIDStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid build ID format")
		return
	}

	executions, err := h.buildExecutionService.GetExecutionsByBuildID(r.Context(), buildID)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to get build executions")
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, executions)
}

// HexagonalFailuresHandler handles HTTP requests for failure operations using hexagonal architecture
type HexagonalFailuresHandler struct {
	failureService domain.FailureService
}

func NewHexagonalFailuresHandler(failureService domain.FailureService) *HexagonalFailuresHandler {
	return &HexagonalFailuresHandler{failureService: failureService}
}

func (h *HexagonalFailuresHandler) GetBuildFailures(w http.ResponseWriter, r *http.Request) {
	// This would need to be implemented based on business requirements
	// For now, return empty array as placeholder
	utils.RespondWithJSON(w, http.StatusOK, []interface{}{})
}

// HexagonalProjectHandler handles HTTP requests for project operations using hexagonal architecture
type HexagonalProjectHandler struct {
	projectService domain.ProjectService
}

func NewHexagonalProjectHandler(projectService domain.ProjectService) *HexagonalProjectHandler {
	return &HexagonalProjectHandler{projectService: projectService}
}

func (h *HexagonalProjectHandler) GetAllProjectsHexagonal(w http.ResponseWriter, r *http.Request) {
	projects, err := h.projectService.GetAllProjects(r.Context())
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "failed to retrieve projects")
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, projects)
}

func (h *HexagonalProjectHandler) CreateProjectHexagonal(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	project, err := h.projectService.CreateProject(r.Context(), req.Name)
	if err != nil {
		switch err {
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		case domain.ErrDuplicateProject:
			utils.RespondWithError(w, http.StatusConflict, "project with this name already exists")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to create project")
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusCreated, project)
}

func (h *HexagonalProjectHandler) GetProjectByIDHexagonal(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing project ID")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid project ID format")
		return
	}
	project, err := h.projectService.GetProjectByID(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrProjectNotFound:
			utils.RespondWithError(w, http.StatusNotFound, "project not found")
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, project)
}

func (h *HexagonalProjectHandler) UpdateProjectHexagonal(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing project ID")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid project ID format")
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	project, err := h.projectService.UpdateProject(r.Context(), id, req.Name)
	if err != nil {
		switch err {
		case domain.ErrProjectNotFound:
			utils.RespondWithError(w, http.StatusNotFound, "project not found")
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		case domain.ErrDuplicateProject:
			utils.RespondWithError(w, http.StatusConflict, "project with this name already exists")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to update project")
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, project)
}

func (h *HexagonalProjectHandler) DeleteProjectHexagonal(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "missing project ID")
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "invalid project ID format")
		return
	}
	err = h.projectService.DeleteProject(r.Context(), id)
	if err != nil {
		switch err {
		case domain.ErrProjectNotFound:
			utils.RespondWithError(w, http.StatusNotFound, "project not found")
		case domain.ErrInvalidInput:
			utils.RespondWithError(w, http.StatusBadRequest, "invalid input")
		default:
			utils.RespondWithError(w, http.StatusInternalServerError, "failed to delete project")
		}
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "project deleted successfully"})
}
