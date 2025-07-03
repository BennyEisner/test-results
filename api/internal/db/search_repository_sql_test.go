package db

import (
	"database/sql"
	"testing"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewSearchRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSearchRepository(db)
	if repo == nil {
		t.Fatal("expected repository to be created, got nil")
	}

	if repo.DB != db {
		t.Error("expected repository to use the provided database")
	}
}

// setupMockExpectations sets up mock expectations for a search test
func setupMockExpectations(mock sqlmock.Sqlmock, query string, expected []models.SearchResult) {
	rows := sqlmock.NewRows([]string{"type", "id", "name", "url"})
	for _, result := range expected {
		rows.AddRow(result.Type, result.ID, result.Name, result.URL)
	}

	expectedPattern := "%" + query + "%"
	mock.ExpectQuery("SELECT 'project' as type, p.id, p.name, '/projects/' \\|\\| p.id as url").
		WithArgs(expectedPattern).
		WillReturnRows(rows)
}

// validateSearchResults validates that the search results match the expected results
func validateSearchResults(t *testing.T, results, expected []models.SearchResult) {
	if len(results) != len(expected) {
		t.Errorf("expected %d results, got %d", len(expected), len(results))
		return
	}

	for i, expected := range expected {
		if i >= len(results) {
			t.Errorf("missing result at index %d", i)
			continue
		}
		validateSingleResult(t, i, results[i], expected)
	}
}

// validateSingleResult validates a single search result
func validateSingleResult(t *testing.T, index int, result, expected models.SearchResult) {
	if result.Type != expected.Type {
		t.Errorf("result %d: expected Type %s, got %s", index, expected.Type, result.Type)
	}
	if result.ID != expected.ID {
		t.Errorf("result %d: expected ID %d, got %d", index, expected.ID, result.ID)
	}
	if result.Name != expected.Name {
		t.Errorf("result %d: expected Name %s, got %s", index, expected.Name, result.Name)
	}
	if result.URL != expected.URL {
		t.Errorf("result %d: expected URL %s, got %s", index, expected.URL, result.URL)
	}
}

func TestSearchRepository_Search(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSearchRepository(db)

	tests := []struct {
		name        string
		query       string
		expected    []models.SearchResult
		expectError bool
	}{
		{
			name:  "search for projects",
			query: "test",
			expected: []models.SearchResult{
				{
					Type: "project",
					ID:   1,
					Name: "Test Project",
					URL:  "/projects/1",
				},
				{
					Type: "test_suite",
					ID:   2,
					Name: "Test Suite",
					URL:  "/projects/1/suites/2",
				},
				{
					Type: "build",
					ID:   3,
					Name: "build-123",
					URL:  "/projects/1/suites/2/builds/3",
				},
			},
			expectError: false,
		},
		{
			name:  "search with special characters",
			query: "project & co",
			expected: []models.SearchResult{
				{
					Type: "project",
					ID:   4,
					Name: "Project & Co",
					URL:  "/projects/4",
				},
			},
			expectError: false,
		},
		{
			name:        "empty query",
			query:       "",
			expected:    []models.SearchResult{},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupMockExpectations(mock, tt.query, tt.expected)

			results, err := repo.Search(tt.query)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			validateSearchResults(t, results, tt.expected)
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock expectations not met: %v", err)
	}
}

func TestSearchRepository_SearchWithDatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSearchRepository(db)

	// Simulate a database error
	mock.ExpectQuery("SELECT 'project' as type, p.id, p.name, '/projects/' \\|\\| p.id as url").
		WithArgs("%test%").
		WillReturnError(sql.ErrConnDone)

	results, err := repo.Search("test")

	if err == nil {
		t.Error("expected error but got none")
	}

	if results != nil {
		t.Error("expected nil results on error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock expectations not met: %v", err)
	}
}

func TestSearchRepository_SearchWithScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock: %v", err)
	}
	defer db.Close()

	repo := NewSearchRepository(db)

	// Return rows with invalid data that will cause scan error
	rows := sqlmock.NewRows([]string{"type", "id", "name", "url"}).
		AddRow("project", "invalid_id", "Test Project", "/projects/1")

	mock.ExpectQuery("SELECT 'project' as type, p.id, p.name, '/projects/' \\|\\| p.id as url").
		WithArgs("%test%").
		WillReturnRows(rows)

	results, err := repo.Search("test")

	if err == nil {
		t.Error("expected error but got none")
	}

	if results != nil {
		t.Error("expected nil results on error")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("mock expectations not met: %v", err)
	}
}
