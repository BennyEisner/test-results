package handler

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/BennyEisner/test-results/internal/service"
	"github.com/BennyEisner/test-results/internal/utils"
)

// JUnitImportHandler will process JUnit XML imports.
type JUnitImportHandler struct {
	importService service.JUnitImportServiceInterface
}

// NewJUnitImportHandler creates a JUnitImportHandler.
func NewJUnitImportHandler(is service.JUnitImportServiceInterface) *JUnitImportHandler {
	return &JUnitImportHandler{importService: is}
}

// HandleJUnitImport handles the import of a JUnit XML file for a specific project and suite.
// It expects a POST request to /api/projects/{projectID}/suites/{suiteID}/junit_imports with a multipart form containing a 'junitFile'.
// A new Build will be created by the service layer for this import.
func (jih *JUnitImportHandler) HandleJUnitImport(w http.ResponseWriter, r *http.Request) {
	if err := checkPostMethod(r); err != nil {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, err.Error())
		return
	}

	projectID, suiteID, err := parseProjectAndSuiteIDs(r.URL.Path)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fileBytes, fileName, err := extractJUnitFileFromRequest(r)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	parsedData, err := parseJUnitXML(fileBytes)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdBuild, processingErrors, err := jih.importService.ProcessJUnitData(projectID, suiteID, parsedData)
	if err != nil {
		log.Printf("Error processing JUnit data for project %d, suite %d: %v\n", projectID, suiteID, err)
		if strings.Contains(err.Error(), "not found") {
			utils.RespondWithError(w, http.StatusNotFound, "Error processing JUnit data: "+err.Error())
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error processing JUnit data: "+err.Error())
		}
		return
	}

	respondWithImportResult(w, projectID, suiteID, fileName, createdBuild, processingErrors)
}

// checkPostMethod ensures the request method is POST. Returns an error if not.
func checkPostMethod(r *http.Request) error {
	if r.Method != http.MethodPost {
		return fmt.Errorf("Only POST method is allowed.")
	}
	return nil
}

// parseProjectAndSuiteIDs extracts the project and suite IDs from the request path.
// Returns the IDs or an error if the path is invalid or IDs cannot be parsed.
func parseProjectAndSuiteIDs(path string) (int64, int64, error) {
	pathSegments := strings.Split(strings.TrimPrefix(path, "/api/projects/"), "/")
	if len(pathSegments) != 4 || pathSegments[1] != "suites" || pathSegments[3] != "junit_imports" {
		return 0, 0, fmt.Errorf("Invalid URL. Expected /api/projects/{projectID}/suites/{suiteID}/junit_imports")
	}
	projectID, err := strconv.ParseInt(pathSegments[0], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid project ID: %v", err)
	}
	suiteID, err := strconv.ParseInt(pathSegments[2], 10, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("Invalid suite ID: %v", err)
	}
	return projectID, suiteID, nil
}

// extractJUnitFileFromRequest extracts the JUnit file bytes and filename from the multipart form in the request.
// Returns the file bytes, filename, or an error if extraction fails.
func extractJUnitFileFromRequest(r *http.Request) ([]byte, string, error) {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		return nil, "", fmt.Errorf("Content-Type header must be multipart/form-data")
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, "", fmt.Errorf("Error parsing multipart form: %v", err)
	}
	file, handler, err := r.FormFile("junitFile")
	if err != nil {
		return nil, "", fmt.Errorf("Error retrieving file 'junitFile': %v", err)
	}
	defer file.Close()
	log.Printf("Received JUnit import: File %s, Size %d\n", handler.Filename, handler.Size)
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, "", fmt.Errorf("Error reading file content: %v", err)
	}
	return fileBytes, handler.Filename, nil
}

// parseJUnitXML attempts to unmarshal the provided bytes into JUnitTestSuites.
// Handles both <testsuites> and <testsuite> roots. Returns the parsed data or an error.
func parseJUnitXML(fileBytes []byte) (*models.JUnitTestSuites, error) {
	var parsedData models.JUnitTestSuites
	if err := xml.Unmarshal(fileBytes, &parsedData); err != nil {
		var singleSuite models.JUnitTestSuite
		if errSingle := xml.Unmarshal(fileBytes, &singleSuite); errSingle == nil {
			parsedData.TestSuites = []models.JUnitTestSuite{singleSuite}
			parsedData.Name = singleSuite.Name
			log.Println("Successfully unmarshalled as a single <testsuite> root, wrapped into JUnitTestSuites.")
			return &parsedData, nil
		} else {
			return nil, fmt.Errorf("Invalid XML. Not <testsuites> (%v) or <testsuite> (%v) root.", err, errSingle)
		}
	} else {
		log.Println("Successfully unmarshalled as <testsuites> root.")
	}
	return &parsedData, nil
}

// respondWithImportResult writes the API response for a JUnit import, including build and error info.
func respondWithImportResult(w http.ResponseWriter, projectID, suiteID int64, fileName string, createdBuild *models.Build, processingErrors []string) {
	response := map[string]interface{}{
		"message":          "JUnit XML processed.",
		"projectID":        projectID,
		"suiteID":          suiteID,
		"fileName":         fileName,
		"createdBuildID":   createdBuild.ID,
		"processingErrors": processingErrors,
	}
	if len(processingErrors) > 0 {
		log.Printf("JUnit import for project %d, suite %d, build %d completed with %d processing errors.\n", projectID, suiteID, createdBuild.ID, len(processingErrors))
		utils.RespondWithJSON(w, http.StatusOK, response)
	} else {
		log.Printf("JUnit import for project %d, suite %d, build %d completed successfully.\n", projectID, suiteID, createdBuild.ID)
		utils.RespondWithJSON(w, http.StatusCreated, response)
	}
}
