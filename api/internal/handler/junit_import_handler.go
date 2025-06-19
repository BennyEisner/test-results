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
// A new Build will be created by the service layer for this import.
// Expected path: POST /api/projects/{projectID}/suites/{suiteID}/junit_imports
func (jih *JUnitImportHandler) HandleJUnitImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.RespondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed.")
		return
	}

	pathSegments := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/projects/"), "/")
	// Expected: {projectID}/suites/{suiteID}/junit_imports
	if len(pathSegments) != 4 || pathSegments[1] != "suites" || pathSegments[3] != "junit_imports" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid URL. Expected /api/projects/{projectID}/suites/{suiteID}/junit_imports")
		return
	}

	projectID, err := strconv.ParseInt(pathSegments[0], 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid project ID: "+err.Error())
		return
	}

	suiteID, err := strconv.ParseInt(pathSegments[2], 10, 64)
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid suite ID: "+err.Error())
		return
	}

	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		utils.RespondWithError(w, http.StatusUnsupportedMediaType, "Content-Type header must be multipart/form-data")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil { // Max 10 MB file
		utils.RespondWithError(w, http.StatusBadRequest, "Error parsing multipart form: "+err.Error())
		return
	}

	file, handler, err := r.FormFile("junitFile") // Expected form field name
	if err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Error retrieving file 'junitFile': "+err.Error())
		return
	}
	defer file.Close()

	log.Printf("Received JUnit import for project ID %d, suite ID %d: File %s, Size %d\n", projectID, suiteID, handler.Filename, handler.Size)

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error reading file content: "+err.Error())
		return
	}

	var parsedData models.JUnitTestSuites // Assuming the XML root is <testsuites>
	if err := xml.Unmarshal(fileBytes, &parsedData); err != nil {
		// If root is <testsuite> instead of <testsuites>
		var singleSuite models.JUnitTestSuite
		if errSingle := xml.Unmarshal(fileBytes, &singleSuite); errSingle == nil {
			parsedData.TestSuites = []models.JUnitTestSuite{singleSuite}
			parsedData.Name = singleSuite.Name // Or derive from context
			log.Println("Successfully unmarshalled as a single <testsuite> root, wrapped into JUnitTestSuites.")
		} else {
			utils.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid XML. Not <testsuites> (%v) or <testsuite> (%v) root.", err, errSingle))
			return
		}
	} else {
		log.Println("Successfully unmarshalled as <testsuites> root.")
	}

	// It's crucial that the content of `parsedData` (especially the suite names within it)
	// aligns with the `suiteID` provided in the URL if strict validation is needed.
	// For now, we pass `suiteID` from URL and `parsedData` to the service.
	// The service can perform further validation.

	createdBuild, processingErrors, err := jih.importService.ProcessJUnitData(projectID, suiteID, &parsedData)
	if err != nil {
		// This is a fatal error from the service (e.g., DB transaction failed, suite not found)
		// Log the detailed error server-side
		log.Printf("Error processing JUnit data for project %d, suite %d: %v\n", projectID, suiteID, err)
		// Respond with a generic server error or a more specific one if appropriate
		// (e.g., if err indicates a client-side issue like suite not found, could be 404 or 400)
		if strings.Contains(err.Error(), "not found") { // Basic check, can be more sophisticated
			utils.RespondWithError(w, http.StatusNotFound, "Error processing JUnit data: "+err.Error())
		} else {
			utils.RespondWithError(w, http.StatusInternalServerError, "Error processing JUnit data: "+err.Error())
		}
		return
	}

	// Successfully processed, potentially with non-fatal errors
	response := map[string]interface{}{
		"message":          "JUnit XML processed.",
		"projectID":        projectID,
		"suiteID":          suiteID,
		"fileName":         handler.Filename,
		"createdBuildID":   createdBuild.ID,  // Assuming createdBuild is not nil if err is nil
		"processingErrors": processingErrors, // Will be an empty list if no errors
	}

	if len(processingErrors) > 0 {
		log.Printf("JUnit import for project %d, suite %d, build %d completed with %d processing errors.\n", projectID, suiteID, createdBuild.ID, len(processingErrors))
		// Still a success (200 OK or 201 Created), but with error details
		utils.RespondWithJSON(w, http.StatusOK, response) // Or http.StatusCreated if that's more appropriate
	} else {
		log.Printf("JUnit import for project %d, suite %d, build %d completed successfully.\n", projectID, suiteID, createdBuild.ID)
		utils.RespondWithJSON(w, http.StatusCreated, response) // 201 Created for successful resource creation
	}
}
