package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/BennyEisner/test-results/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSearchService is a mock implementation of domain.SearchService
type MockSearchService struct {
	mock.Mock
}

func (m *MockSearchService) Search(ctx context.Context, query string) ([]*domain.SearchResult, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.SearchResult), args.Error(1)
}

func TestSearchHandler_HandleSearch(t *testing.T) {
	tests := []struct {
		name           string
		query          string
		mockResults    []*domain.SearchResult
		mockError      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:  "successful search",
			query: "test",
			mockResults: []*domain.SearchResult{
				{Type: "project", ID: 1, Name: "Test Project", URL: "/projects/1"},
				{Type: "test_suite", ID: 2, Name: "Test Suite", URL: "/projects/1/suites/2"},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody: []*domain.SearchResult{
				{Type: "project", ID: 1, Name: "Test Project", URL: "/projects/1"},
				{Type: "test_suite", ID: 2, Name: "Test Suite", URL: "/projects/1/suites/2"},
			},
		},
		{
			name:           "empty query",
			query:          "",
			mockResults:    []*domain.SearchResult{},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedBody:   []*domain.SearchResult{},
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

			// Set up mock expectations
			mockService.On("Search", mock.Anything, tt.query).Return(tt.mockResults, tt.mockError)

			// Call handler
			handler.HandleSearch(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var results []*domain.SearchResult
				err := json.Unmarshal(w.Body.Bytes(), &results)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, results)
			} else {
				var errorResponse map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &errorResponse)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedBody, errorResponse)
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSearchHandler_HandleSearch_NoQuery(t *testing.T) {
	mockService := new(MockSearchService)
	handler := NewSearchHandler(mockService)

	// Create request without query parameter
	req := httptest.NewRequest("GET", "/api/search", nil)
	w := httptest.NewRecorder()

	// Set up mock expectations
	mockService.On("Search", mock.Anything, "").Return([]*domain.SearchResult{}, nil)

	// Call handler
	handler.HandleSearch(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)

	var results []*domain.SearchResult
	err := json.Unmarshal(w.Body.Bytes(), &results)
	assert.NoError(t, err)
	assert.Equal(t, []*domain.SearchResult{}, results)

	mockService.AssertExpectations(t)
}
