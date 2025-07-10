package application

import (
	"context"
	"testing"

	"github.com/BennyEisner/test-results/internal/project/domain"
	"github.com/BennyEisner/test-results/internal/test_case/application"
	"github.com/BennyEisner/test-results/internal/test_case/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTestCaseRepository is a mock implementation of TestCaseRepository
type MockTestCaseRepository struct {
	mock.Mock
}

func (m *MockTestCaseRepository) GetByID(ctx context.Context, id int64) (*models.TestCase, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) GetAllBySuiteID(ctx context.Context, suiteID int64) ([]*models.TestCase, error) {
	args := m.Called(ctx, suiteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) GetByName(ctx context.Context, suiteID int64, name string) (*models.TestCase, error) {
	args := m.Called(ctx, suiteID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) Create(ctx context.Context, tc *models.TestCase) error {
	args := m.Called(ctx, tc)
	return args.Error(0)
}

func (m *MockTestCaseRepository) Update(ctx context.Context, id int64, name, classname string) (*models.TestCase, error) {
	args := m.Called(ctx, id, name, classname)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTestCaseService_GetTestCaseByID(t *testing.T) {
	mockRepo := new(MockTestCaseRepository)
	service := application.NewTestCaseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedTestCase := &models.TestCase{
			ID:        1,
			SuiteID:   123,
			Name:      "TestExample",
			Classname: "com.example.TestExample",
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(expectedTestCase, nil).Once()

		result, err := service.GetTestCase(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedTestCase, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.GetTestCase(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidTestCaseName, err)
		assert.Nil(t, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, nil).Once()

		result, err := service.GetTestCase(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrTestCaseNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestCaseService_GetTestCasesBySuiteID(t *testing.T) {
	mockRepo := new(MockTestCaseRepository)
	service := application.NewTestCaseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedTestCases := []*models.TestCase{
			{
				ID:        1,
				SuiteID:   123,
				Name:      "TestExample1",
				Classname: "com.example.TestExample1",
			},
			{
				ID:        2,
				SuiteID:   123,
				Name:      "TestExample2",
				Classname: "com.example.TestExample2",
			},
		}

		mockRepo.On("GetAllBySuiteID", ctx, int64(123)).Return(expectedTestCases, nil).Once()

		result, err := service.GetTestCasesBySuite(ctx, 123)

		assert.NoError(t, err)
		assert.Equal(t, expectedTestCases, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.GetTestCasesBySuite(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidTestSuiteName, err)
		assert.Nil(t, result)
	})
}

func TestTestCaseService_CreateTestCase(t *testing.T) {
	mockRepo := new(MockTestCaseRepository)
	service := application.NewTestCaseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		suiteID := int64(123)
		name := "TestExample"
		classname := "com.example.TestExample"

		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.TestCase")).Return(nil).Once()

		result, err := service.CreateTestCase(ctx, suiteID, name, classname)

		assert.NoError(t, err)
		assert.Equal(t, suiteID, result.SuiteID)
		assert.Equal(t, name, result.Name)
		assert.Equal(t, classname, result.Classname)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input - empty name", func(t *testing.T) {
		result, err := service.CreateTestCase(ctx, 123, "", "classname")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidTestCaseName, err)
		assert.Nil(t, result)
	})

	t.Run("invalid input - empty classname", func(t *testing.T) {
		result, err := service.CreateTestCase(ctx, 123, "name", "")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidTestCaseName, err)
		assert.Nil(t, result)
	})

}

func TestTestCaseService_UpdateTestCase(t *testing.T) {
	mockRepo := new(MockTestCaseRepository)
	service := application.NewTestCaseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		testCaseID := int64(1)
		name := "UpdatedTest"
		classname := "com.example.UpdatedTest"
		updatedTestCase := &models.TestCase{
			ID:        testCaseID,
			SuiteID:   123,
			Name:      name,
			Classname: classname,
		}

		mockRepo.On("Update", ctx, testCaseID, name, classname).Return(updatedTestCase, nil).Once()

		result, err := service.UpdateTestCase(ctx, testCaseID, name, classname)

		assert.NoError(t, err)
		assert.Equal(t, updatedTestCase, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.UpdateTestCase(ctx, 0, "name", "classname")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidTestCaseName, err)
		assert.Nil(t, result)
	})

}

func TestTestCaseService_DeleteTestCase(t *testing.T) {
	mockRepo := new(MockTestCaseRepository)
	service := application.NewTestCaseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		testCaseID := int64(1)

		mockRepo.On("Delete", ctx, testCaseID).Return(nil).Once()

		err := service.DeleteTestCase(ctx, testCaseID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		err := service.DeleteTestCase(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidTestCaseName, err)
	})

}
