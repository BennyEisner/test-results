package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSearchService is a mock implementation of ports.SearchService
type MockSearchService struct {
	mock.Mock
}

func (m *MockSearchService) Search(ctx context.Context, query string) ([]*models.SearchResult, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SearchResult), args.Error(1)
}

func TestSearchHandler_Search(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockResults    []*models.SearchResult
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:  "successful search",
			query: "test",
			mockResults: []*models.SearchResult{
				{Type: "project", ID: 1, Name: "Test Project", URL: "/projects/1"},
				{Type: "test_suite", ID: 2, Name: "Test Suite", URL: "/projects/1/suites/2"},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: []*models.SearchResult{
				{Type: "project", ID: 1, Name: "Test Project", URL: "/projects/1"},
				{Type: "test_suite", ID: 2, Name: "Test Suite", URL: "/projects/1/suites/2"},
			},
		},
		{
			name:           "empty query",
			query:          "",
			mockResults:    []*models.SearchResult{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "missing query parameter"},
		},
		{
			name:           "service error",
			query:          "test",
			mockResults:    nil,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   map[string]string{"error": "failed to perform search"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSearchService)
			handler := NewSearchHandler(mockService)

			// Create request
			req := httptest.NewRequest("GET", "/api/search?q="+tt.query, nil)
			w := httptest.NewRecorder()

			// Set up mock expectations (only for non-empty queries)
			if tt.query != "" {
				mockService.On("Search", mock.Anything, tt.query).Return(tt.mockResults, tt.mockError)
			}

			// Call handler
			handler.Search(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var results []*models.SearchResult
				err := json.Unmarshal(w.Body.Bytes(), &results)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, results)
			} else {
				// For error responses, check the body contains the error message
				body := w.Body.String()
				if tt.expectedStatus == http.StatusBadRequest {
					assert.Contains(t, body, "missing query parameter")
				} else if tt.expectedStatus == http.StatusInternalServerError {
					assert.Contains(t, body, "database error")
				}
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSearchHandler_Search_NoQuery(t *testing.T) {
	mockService := new(MockSearchService)
	handler := NewSearchHandler(mockService)

	// Create request without query parameter
	req := httptest.NewRequest("GET", "/api/search", nil)
	w := httptest.NewRecorder()

	// Call handler
	handler.Search(w, req)

	// Assert response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Check error message
	body := w.Body.String()
	assert.Contains(t, body, "missing query parameter")

	// No mock expectations since the service is not called for empty queries
}
