package client

import (
	"bytes"
	"fmt"
	"github.com/BennyEisner/test-results/cli/internal/config"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// APIClient makes requests to the test-results API
type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewAPICLient creates a new API client
func newAPIClient(cfg *config.Config) *APIClient {
	return &APIClient{
		BaseURL:    cfg.APIBaseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *APIClient) PostJUnitFile(projectID, suiteID int64, filePath string) (string, error) {
	url := fmt.Sprintf("%s/api/projects/%d/suites/%d/junit_imports", c.BaseURL, projectID, suiteID)

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("Error openening file: %w", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("junitFile", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("Error creating form file: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("Error copying file content: %w", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return "", fmt.Errorf("Error creating request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending request: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body: %w", err)
	}
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("API error (status: %d): %s", resp.StatusCode, string(respBody))
	}
	return string(respBody), nil
}
