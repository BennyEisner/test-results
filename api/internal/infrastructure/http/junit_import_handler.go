package http

import (
	"encoding/xml"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/BennyEisner/test-results/internal/domain"
)

type JUnitImportHandler struct {
	importService domain.JUnitImportService
}

func NewJUnitImportHandler(importService domain.JUnitImportService) *JUnitImportHandler {
	return &JUnitImportHandler{importService: importService}
}

// HandleJUnitImport handles the import of a JUnit XML file for a specific project and suite.
// It expects a POST request to /api/projects/{projectID}/suites/{suiteID}/junit_imports with a multipart form containing a 'junitFile'.
func (h *JUnitImportHandler) HandleJUnitImport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondWithError(w, http.StatusMethodNotAllowed, "Only POST method is allowed.")
		return
	}

	projectID, suiteID, err := parseProjectAndSuiteIDs(r.URL.Path)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fileBytes, fileName, err := extractJUnitFileFromRequest(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	parsedData, err := parseJUnitXML(fileBytes)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdBuild, err := h.importService.ProcessJUnitData(r.Context(), projectID, suiteID, parsedData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error processing JUnit data: "+err.Error())
		return
	}

	response := map[string]interface{}{
		"message":        "JUnit XML processed.",
		"projectID":      projectID,
		"suiteID":        suiteID,
		"fileName":       fileName,
		"createdBuildID": createdBuild.ID,
	}
	respondWithJSON(w, http.StatusCreated, response)
}

func parseProjectAndSuiteIDs(path string) (int64, int64, error) {
	pathSegments := strings.Split(strings.TrimPrefix(path, "/api/projects/"), "/")
	if len(pathSegments) != 4 || pathSegments[1] != "suites" || pathSegments[3] != "junit_imports" {
		return 0, 0, http.ErrNotSupported
	}
	projectID, err := strconv.ParseInt(pathSegments[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	suiteID, err := strconv.ParseInt(pathSegments[2], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	return projectID, suiteID, nil
}

func extractJUnitFileFromRequest(r *http.Request) ([]byte, string, error) {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		return nil, "", http.ErrNotSupported
	}
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		return nil, "", err
	}
	file, handler, err := r.FormFile("junitFile")
	if err != nil {
		return nil, "", err
	}
	defer file.Close()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, "", err
	}
	return fileBytes, handler.Filename, nil
}

func parseJUnitXML(fileBytes []byte) (*domain.JUnitTestSuites, error) {
	var parsedData domain.JUnitTestSuites
	if err := xml.Unmarshal(fileBytes, &parsedData); err != nil {
		var singleSuite domain.JUnitTestSuite
		if errSingle := xml.Unmarshal(fileBytes, &singleSuite); errSingle == nil {
			parsedData.TestSuites = []domain.JUnitTestSuite{singleSuite}
			parsedData.Name = singleSuite.Name
			return &parsedData, nil
		} else {
			return nil, err
		}
	}
	return &parsedData, nil
}
