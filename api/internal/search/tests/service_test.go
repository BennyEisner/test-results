package application

import (
	"context"
	"errors"
	"testing"

	"github.com/BennyEisner/test-results/internal/search/application"
	"github.com/BennyEisner/test-results/internal/search/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSearchRepository is a mock implementation of domain.SearchRepository
type MockSearchRepository struct {
	mock.Mock
}

func (m *MockSearchRepository) Search(ctx context.Context, query string) ([]*models.SearchResult, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.SearchResult), args.Error(1)
}

func TestSearchService_Search(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		mockResults     []*models.SearchResult
		mockError       error
		expectedResults []*models.SearchResult
		expectedError   error
	}{
		{
			name:  "successful search",
			query: "test",
			mockResults: []*models.SearchResult{
				{Type: "project", ID: 1, Name: "Test Project", URL: "/projects/1"},
				{Type: "test_suite", ID: 2, Name: "Test Suite", URL: "/projects/1/suites/2"},
			},
			mockError: nil,
			expectedResults: []*models.SearchResult{
				{Type: "project", ID: 1, Name: "Test Project", URL: "/projects/1"},
				{Type: "test_suite", ID: 2, Name: "Test Suite", URL: "/projects/1/suites/2"},
			},
			expectedError: nil,
		},
		{
			name:            "empty query returns empty results",
			query:           "",
			mockResults:     nil,
			mockError:       nil,
			expectedResults: []*models.SearchResult{},
			expectedError:   nil,
		},
		{
			name:            "repository error",
			query:           "test",
			mockResults:     nil,
			mockError:       errors.New("database error"),
			expectedResults: nil,
			expectedError:   errors.New("failed to search: database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSearchRepository)
			service := application.NewSearchService(mockRepo)

			if tt.query != "" {
				mockRepo.On("Search", mock.Anything, tt.query).Return(tt.mockResults, tt.mockError)
			}

			results, err := service.Search(context.Background(), tt.query)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResults, results)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
