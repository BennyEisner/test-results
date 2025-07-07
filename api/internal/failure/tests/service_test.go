package application

import (
	"context"
	"testing"

	"github.com/BennyEisner/test-results/internal/domain"
	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockFailureRepository is a mock implementation of FailureRepository
type MockFailureRepository struct {
	mock.Mock
}

func (m *MockFailureRepository) GetByID(ctx context.Context, id int64) (*models.Failure, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Failure), args.Error(1)
}

func (m *MockFailureRepository) GetByExecutionID(ctx context.Context, executionID int64) (*models.Failure, error) {
	args := m.Called(ctx, executionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Failure), args.Error(1)
}

func (m *MockFailureRepository) Create(ctx context.Context, failure *models.Failure) error {
	args := m.Called(ctx, failure)
	return args.Error(0)
}

func (m *MockFailureRepository) Update(ctx context.Context, id int64, failure *models.Failure) (*models.Failure, error) {
	args := m.Called(ctx, id, failure)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Failure), args.Error(1)
}

func (m *MockFailureRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestFailureService_GetFailureByID(t *testing.T) {
	mockRepo := new(MockFailureRepository)
	service := NewFailureService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedFailure := &models.Failure{
			ID:                       1,
			BuildTestCaseExecutionID: 123,
			Message:                  stringPtr("Test failure"),
			Type:                     stringPtr("AssertionError"),
			Details:                  stringPtr("Expected true but got false"),
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(expectedFailure, nil).Once()

		result, err := service.GetFailureByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedFailure, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.GetFailureByID(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, "invalid failure ID", err.Error())
		assert.Nil(t, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, nil).Once()

		result, err := service.GetFailureByID(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrFailureNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestFailureService_CreateFailure(t *testing.T) {
	mockRepo := new(MockFailureRepository)
	service := NewFailureService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		failure := &models.Failure{
			BuildTestCaseExecutionID: 123,
			Message:                  stringPtr("Test failure"),
			Type:                     stringPtr("AssertionError"),
			Details:                  stringPtr("Expected true but got false"),
		}

		mockRepo.On("Create", ctx, failure).Return(nil).Once()

		result, err := service.CreateFailure(ctx, failure)

		assert.NoError(t, err)
		assert.Equal(t, failure, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input - nil failure", func(t *testing.T) {
		result, err := service.CreateFailure(ctx, nil)

		assert.Error(t, err)
		assert.Equal(t, "failure cannot be nil", err.Error())
		assert.Nil(t, result)
	})

	t.Run("invalid input - invalid execution ID", func(t *testing.T) {
		failure := &models.Failure{
			BuildTestCaseExecutionID: 0,
		}

		result, err := service.CreateFailure(ctx, failure)

		assert.Error(t, err)
		assert.Equal(t, "execution ID is required", err.Error())
		assert.Nil(t, result)
	})
}

func TestFailureService_DeleteFailure(t *testing.T) {
	mockRepo := new(MockFailureRepository)
	service := NewFailureService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Delete", ctx, int64(1)).Return(nil).Once()

		err := service.DeleteFailure(ctx, 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		err := service.DeleteFailure(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, "invalid failure ID", err.Error())
	})
}

// Helper function to create string pointers
func stringPtr(s string) *string {
	return &s
}
