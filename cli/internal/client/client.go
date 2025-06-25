package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/BennyEisner/test-results/cli/internal/config"
)

// APIClient makes requests to the test-results API
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAPICLient creates a new API client
func NewAPIClient(cfg *config.Config) *APIClient {
	return &APIClient{
		BaseURL:    cfg.APIBaseURL,
		HTTPClient: &http.Client{},
	}
}

// PostJUnitFile uploads a JUnit XML file to the API.
func (c *APIClient) PostJUnitFile(projectID, suiteID int64, filePath string) (string, error) {
	url := fmt.Sprintf("%s/api/projects/%d/suites/%d/junit_imports", c.BaseURL, projectID, suiteID)

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	// Create a buffer and multipart writer
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// Create a form file field
	fileWriter, err := writer.CreateFormFile("junitFile", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("error creating form file: %w", err)
	}

	_, err = io.Copy(fileWriter, file)
	if err != nil {
		return "", fmt.Errorf("error copying file content: %w", err)
	}

	// Close the writer before creating the request
	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("error closing multipart writer: %w", err)
	}

	//  Create the HTTP request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	//  Set headers
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Length", fmt.Sprintf("%d", requestBody.Len()))

	//  Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	//  Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	return string(respBody), nil
}
