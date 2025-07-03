package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/BennyEisner/test-results/internal/utils"
)

// mockProjectService implements ProjectServiceInterface for testing
// Each method can be set to a custom function in tests
// If not set, it panics

type mockProjectService struct {
	GetProjectByIDFunc        func(id int64) (*models.Project, error)
	GetAllProjectsFunc        func() ([]models.Project, error)
	CreateProjectFunc         func(name string) (*models.Project, error)
	UpdateProjectFunc         func(id int64, name string) (*models.Project, error)
	DeleteProjectFunc         func(id int64) (int64, error)
	GetDBTestProjectCountFunc func() (int, error)
	GetProjectByNameFunc      func(name string) (*models.Project, error)
}

func (m *mockProjectService) GetProjectByID(id int64) (*models.Project, error) {
	return m.GetProjectByIDFunc(id)
}
func (m *mockProjectService) GetAllProjects() ([]models.Project, error) {
	return m.GetAllProjectsFunc()
}
func (m *mockProjectService) CreateProject(name string) (*models.Project, error) {
	return m.CreateProjectFunc(name)
}
func (m *mockProjectService) UpdateProject(id int64, name string) (*models.Project, error) {
	return m.UpdateProjectFunc(id, name)
}
func (m *mockProjectService) DeleteProject(id int64) (int64, error) {
	return m.DeleteProjectFunc(id)
}
func (m *mockProjectService) GetDBTestProjectCount() (int, error) {
	return m.GetDBTestProjectCountFunc()
}
func (m *mockProjectService) GetProjectByName(name string) (*models.Project, error) {
	return m.GetProjectByNameFunc(name)
}

func TestProjectHandler_GetProjectByID(t *testing.T) {
	mockSvc := &mockProjectService{
		GetProjectByIDFunc: func(id int64) (*models.Project, error) {
			if id == 1 {
				return &models.Project{ID: 1, Name: "Test Project"}, nil
			}
			if id == 404 {
				return nil, sql.ErrNoRows
			}
			return nil, errors.New("db error")
		},
	}
	h := NewProjectHandler(mockSvc)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{"valid", "1", http.StatusOK, "Test Project"},
		{"not found", "404", http.StatusNotFound, "Project not found"},
		{"db error", "500", http.StatusInternalServerError, "Database error"},
		{"bad id", "abc", http.StatusBadRequest, "Invalid project ID format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/api/projects/"+tt.id, nil)
			req.SetPathValue("id", tt.id)
			rr := httptest.NewRecorder()

			h.GetProjectByID(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
			if !strings.Contains(rr.Body.String(), tt.expectedBody) {
				t.Errorf("expected body to contain %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

// setupGetProjectsTest sets up the test environment for GetProjects tests
func setupGetProjectsTest() *ProjectHandler {
	mockSvc := &mockProjectService{
		GetAllProjectsFunc: func() ([]models.Project, error) {
			return []models.Project{{ID: 1, Name: "A"}, {ID: 2, Name: "B"}}, nil
		},
		GetProjectByNameFunc: func(name string) (*models.Project, error) {
			if name == "A" {
				return &models.Project{ID: 1, Name: "A"}, nil
			}
			return nil, sql.ErrNoRows
		},
	}
	return NewProjectHandler(mockSvc)
}

// validateProjectsResponse validates the response for projects endpoint
func validateProjectsResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int, expectedCount int) {
	if rr.Code != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, rr.Code)
	}

	if expectedCount > 0 {
		var got []utils.Project
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if len(got) != expectedCount {
			t.Errorf("expected %d projects, got %d", expectedCount, len(got))
		}
	}
}

// validateSingleProjectResponse validates the response for a single project
func validateSingleProjectResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int, expectedName string) {
	if rr.Code != expectedStatus {
		t.Errorf("expected status %d, got %d", expectedStatus, rr.Code)
	}

	if expectedName != "" {
		var got utils.Project
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.Name != expectedName {
			t.Errorf("expected project name %s, got %s", expectedName, got.Name)
		}
	}
}

func TestProjectHandler_GetProjects(t *testing.T) {
	h := setupGetProjectsTest()

	t.Run("all projects", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/projects", nil)
		rr := httptest.NewRecorder()
		h.GetProjects(rr, req)
		validateProjectsResponse(t, rr, http.StatusOK, 2)
	})

	t.Run("by name found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/projects?name=A", nil)
		rr := httptest.NewRecorder()
		h.GetProjects(rr, req)
		validateSingleProjectResponse(t, rr, http.StatusOK, "A")
	})

	t.Run("by name not found", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/projects?name=Z", nil)
		rr := httptest.NewRecorder()
		h.GetProjects(rr, req)
		if rr.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Project not found") {
			t.Errorf("expected not found message, got %q", rr.Body.String())
		}
	})
}

func TestProjectHandler_CreateProject_JSON(t *testing.T) {
	mockSvc := &mockProjectService{
		CreateProjectFunc: func(name string) (*models.Project, error) {
			if name == "" {
				return nil, errors.New("empty name")
			}
			return &models.Project{ID: 1, Name: name}, nil
		},
	}
	h := NewProjectHandler(mockSvc)

	t.Run("valid json", func(t *testing.T) {
		body := `{"name":"Test Project"}`
		req := httptest.NewRequest("POST", "/api/projects", io.NopCloser(strings.NewReader(body)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateProject(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", rr.Code)
		}
		var got utils.Project
		if err := json.NewDecoder(rr.Body).Decode(&got); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}
		if got.Name != "Test Project" {
			t.Errorf("expected project name 'Test Project', got %q", got.Name)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		body := `{"name":123}`
		req := httptest.NewRequest("POST", "/api/projects", io.NopCloser(strings.NewReader(body)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateProject(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Invalid JSON") {
			t.Errorf("expected invalid JSON error, got %q", rr.Body.String())
		}
	})

	t.Run("empty name", func(t *testing.T) {
		body := `{"name":""}`
		req := httptest.NewRequest("POST", "/api/projects", io.NopCloser(strings.NewReader(body)))
		req.Header.Set("Content-Type", "application/json")
		rr := httptest.NewRecorder()
		h.CreateProject(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "Project Name is required") {
			t.Errorf("expected required name error, got %q", rr.Body.String())
		}
	})
}
