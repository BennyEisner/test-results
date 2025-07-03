package models

import (
	"database/sql/driver"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestProject_ScanFromRow(t *testing.T) {
	tests := []struct {
		testName string
		id       int64
		name     string
		expected Project
	}{
		{
			testName: "valid project",
			id:       1,
			name:     "Test Project",
			expected: Project{
				ID:   1,
				Name: "Test Project",
			},
		},
		{
			testName: "project with special characters",
			id:       2,
			name:     "Project & Co. (Ltd.)",
			expected: Project{
				ID:   2,
				Name: "Project & Co. (Ltd.)",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(tt.id, tt.name)
			mock.ExpectQuery("SELECT").WillReturnRows(rows)

			row := db.QueryRow("SELECT id, name FROM projects WHERE id = ?", tt.id)
			project := &Project{}
			err = project.ScanFromRow(row)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if project.ID != tt.expected.ID {
				t.Errorf("expected ID %d, got %d", tt.expected.ID, project.ID)
			}

			if project.Name != tt.expected.Name {
				t.Errorf("expected Name %s, got %s", tt.expected.Name, project.Name)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("mock expectations not met: %v", err)
			}
		})
	}
}

// setupMockRows sets up mock rows for the test
func setupMockRows(mock sqlmock.Sqlmock, rows [][]driver.Value) {
	mockRows := sqlmock.NewRows([]string{"id", "name"})
	for _, row := range rows {
		mockRows.AddRow(row...)
	}
	mock.ExpectQuery("SELECT").WillReturnRows(mockRows)
}

// validateProjects validates that the scanned projects match the expected projects
func validateProjects(t *testing.T, projects, expected []Project) {
	if len(projects) != len(expected) {
		t.Errorf("expected %d projects, got %d", len(expected), len(projects))
		return
	}

	for i, expected := range expected {
		if i >= len(projects) {
			t.Errorf("missing project at index %d", i)
			continue
		}
		validateSingleProject(t, i, projects[i], expected)
	}
}

// validateSingleProject validates a single project
func validateSingleProject(t *testing.T, index int, project, expected Project) {
	if project.ID != expected.ID {
		t.Errorf("project %d: expected ID %d, got %d", index, expected.ID, project.ID)
	}
	if project.Name != expected.Name {
		t.Errorf("project %d: expected Name %s, got %s", index, expected.Name, project.Name)
	}
}

func TestProject_ScanFromRows(t *testing.T) {
	tests := []struct {
		testName string
		rows     [][]driver.Value
		expected []Project
	}{
		{
			testName: "single row",
			rows: [][]driver.Value{
				{int64(1), "Project 1"},
			},
			expected: []Project{
				{ID: 1, Name: "Project 1"},
			},
		},
		{
			testName: "multiple rows",
			rows: [][]driver.Value{
				{int64(1), "Project 1"},
				{int64(2), "Project 2"},
				{int64(3), "Project 3"},
			},
			expected: []Project{
				{ID: 1, Name: "Project 1"},
				{ID: 2, Name: "Project 2"},
				{ID: 3, Name: "Project 3"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			setupMockRows(mock, tt.rows)

			rowsResult, err := db.Query("SELECT id, name FROM projects")
			if err != nil {
				t.Fatalf("failed to query: %v", err)
			}
			defer rowsResult.Close()

			var projects []Project
			for rowsResult.Next() {
				project := &Project{}
				err = project.ScanFromRows(rowsResult)
				if err != nil {
					t.Errorf("unexpected error scanning row: %v", err)
				}
				projects = append(projects, *project)
			}

			validateProjects(t, projects, tt.expected)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("mock expectations not met: %v", err)
			}
		})
	}
}

func TestProject_Insert(t *testing.T) {
	tests := []struct {
		testName    string
		projectName string
		expectedID  int64
		expectError bool
	}{
		{
			testName:    "valid project",
			projectName: "New Project",
			expectedID:  1,
			expectError: false,
		},
		{
			testName:    "project with special characters",
			projectName: "Project & Co. (Ltd.)",
			expectedID:  2,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock: %v", err)
			}
			defer db.Close()

			mock.ExpectQuery("INSERT INTO projects").
				WithArgs(tt.projectName).
				WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tt.expectedID))

			project := &Project{Name: tt.projectName}
			err = project.Insert(db)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if project.ID != tt.expectedID {
				t.Errorf("expected ID %d, got %d", tt.expectedID, project.ID)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("mock expectations not met: %v", err)
			}
		})
	}
}

func TestProjectXML(t *testing.T) {
	project := Project{
		ID:   1,
		Name: "Test Project",
	}

	projectXML := ProjectXML{
		Project: project,
	}

	// Test that the struct can be created without errors
	if projectXML.Project.ID != project.ID {
		t.Errorf("expected ID %d, got %d", project.ID, projectXML.Project.ID)
	}

	if projectXML.Project.Name != project.Name {
		t.Errorf("expected Name %s, got %s", project.Name, projectXML.Project.Name)
	}
}
