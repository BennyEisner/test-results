package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewSQLProjectRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSQLProjectRepository(db)
	if repo == nil {
		t.Fatal("expected repository to be created, got nil")
	}

	if repo.DB != db {
		t.Error("expected repository to use the provided database")
	}
}

func TestSQLProjectRepository_GetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSQLProjectRepository(db)

	tests := []struct {
		name        string
		id          int64
		expected    *models.Project
		expectError bool
	}{
		{
			name: "valid project",
			id:   1,
			expected: &models.Project{
				ID:   1,
				Name: "Test Project",
			},
			expectError: false,
		},
		{
			name:        "project not found",
			id:          999,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectError {
				mock.ExpectQuery("SELECT id, name FROM projects WHERE id = \\$1").
					WithArgs(tt.id).
					WillReturnError(sql.ErrNoRows)
			} else {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(tt.expected.ID, tt.expected.Name)
				mock.ExpectQuery("SELECT id, name FROM projects WHERE id = \\$1").
					WithArgs(tt.id).
					WillReturnRows(rows)
			}

			result, err := repo.GetByID(context.Background(), tt.id)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if result.ID != tt.expected.ID {
				t.Errorf("expected ID %d, got %d", tt.expected.ID, result.ID)
			}

			if result.Name != tt.expected.Name {
				t.Errorf("expected Name %s, got %s", tt.expected.Name, result.Name)
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock expectations not met: %v", err)
	}
}

func TestSQLProjectRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSQLProjectRepository(db)

	expectedProjects := []*models.Project{
		{ID: 1, Name: "Project 1"},
		{ID: 2, Name: "Project 2"},
		{ID: 3, Name: "Project 3"},
	}

	rows := sqlmock.NewRows([]string{"id", "name"})
	for _, p := range expectedProjects {
		rows.AddRow(p.ID, p.Name)
	}

	mock.ExpectQuery("SELECT id, name FROM projects").WillReturnRows(rows)

	results, err := repo.GetAll(context.Background())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(results) != len(expectedProjects) {
		t.Errorf("expected %d projects, got %d", len(expectedProjects), len(results))
	}

	for i, expected := range expectedProjects {
		if results[i].ID != expected.ID {
			t.Errorf("project %d: expected ID %d, got %d", i, expected.ID, results[i].ID)
		}
		if results[i].Name != expected.Name {
			t.Errorf("project %d: expected Name %s, got %s", i, expected.Name, results[i].Name)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock expectations not met: %v", err)
	}
}

func TestSQLProjectRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSQLProjectRepository(db)

	project := &models.Project{Name: "New Project"}
	expectedID := int64(1)

	mock.ExpectQuery("INSERT INTO projects\\(name\\) VALUES\\(\\$1\\) RETURNING id").
		WithArgs(project.Name).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))

	err = repo.Create(context.Background(), project)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if project.ID != expectedID {
		t.Errorf("expected ID %d, got %d", expectedID, project.ID)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock expectations not met: %v", err)
	}
}

func TestSQLProjectRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSQLProjectRepository(db)

	projectID := int64(1)

	mock.ExpectExec("DELETE FROM projects WHERE id = \\$1").
		WithArgs(projectID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.Delete(context.Background(), projectID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock expectations not met: %v", err)
	}
}

func TestSQLProjectRepository_GetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSQLProjectRepository(db)

	expectedProject := &models.Project{
		ID:   1,
		Name: "Test Project",
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(expectedProject.ID, expectedProject.Name)

	mock.ExpectQuery("SELECT id, name FROM projects WHERE name = \\$1").
		WithArgs(expectedProject.Name).
		WillReturnRows(rows)

	result, err := repo.GetByName(context.Background(), expectedProject.Name)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.ID != expectedProject.ID {
		t.Errorf("expected ID %d, got %d", expectedProject.ID, result.ID)
	}

	if result.Name != expectedProject.Name {
		t.Errorf("expected Name %s, got %s", expectedProject.Name, result.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock expectations not met: %v", err)
	}
}

func TestSQLProjectRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSQLProjectRepository(db)

	projectID := int64(1)
	newName := "Updated Project"
	expectedProject := &models.Project{
		ID:   projectID,
		Name: newName,
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(expectedProject.ID, expectedProject.Name)

	mock.ExpectQuery("UPDATE projects SET name = \\$1 WHERE id = \\$2 RETURNING id, name").
		WithArgs(newName, projectID).
		WillReturnRows(rows)

	result, err := repo.Update(context.Background(), projectID, newName)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result.ID != expectedProject.ID {
		t.Errorf("expected ID %d, got %d", expectedProject.ID, result.ID)
	}

	if result.Name != expectedProject.Name {
		t.Errorf("expected Name %s, got %s", expectedProject.Name, result.Name)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock expectations not met: %v", err)
	}
}

func TestSQLProjectRepository_Count(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSQLProjectRepository(db)

	expectedCount := 5

	rows := sqlmock.NewRows([]string{"count"}).AddRow(expectedCount)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM projects").WillReturnRows(rows)

	result, err := repo.Count(context.Background())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != expectedCount {
		t.Errorf("expected count %d, got %d", expectedCount, result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock expectations not met: %v", err)
	}
}
