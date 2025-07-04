package application

import (
	"context"
	"testing"

	"github.com/BennyEisner/test-results/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTestCaseRepository is a mock implementation of TestCaseRepository
type MockTestCaseRepository struct {
	mock.Mock
}

func (m *MockTestCaseRepository) GetByID(ctx context.Context, id int64) (*domain.TestCase, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) GetAllBySuiteID(ctx context.Context, suiteID int64) ([]*domain.TestCase, error) {
	args := m.Called(ctx, suiteID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) GetByName(ctx context.Context, suiteID int64, name string) (*domain.TestCase, error) {
	args := m.Called(ctx, suiteID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) Create(ctx context.Context, tc *domain.TestCase) error {
	args := m.Called(ctx, tc)
	return args.Error(0)
}

func (m *MockTestCaseRepository) Update(ctx context.Context, id int64, name, classname string) (*domain.TestCase, error) {
	args := m.Called(ctx, id, name, classname)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTestCaseService_GetTestCaseByID(t *testing.T) {
	mockRepo := new(MockTestCaseRepository)
	service := NewTestCaseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedTestCase := &domain.TestCase{
			ID:        1,
			SuiteID:   123,
			Name:      "TestExample",
			Classname: "com.example.TestExample",
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(expectedTestCase, nil).Once()

		result, err := service.GetTestCaseByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedTestCase, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.GetTestCaseByID(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidInput, err)
		assert.Nil(t, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, nil).Once()

		result, err := service.GetTestCaseByID(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrTestCaseNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestCaseService_GetTestCasesBySuiteID(t *testing.T) {
	mockRepo := new(MockTestCaseRepository)
	service := NewTestCaseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedTestCases := []*domain.TestCase{
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

		result, err := service.GetTestCasesBySuiteID(ctx, 123)

		assert.NoError(t, err)
		assert.Equal(t, expectedTestCases, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.GetTestCasesBySuiteID(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidInput, err)
		assert.Nil(t, result)
	})
}

func TestTestCaseService_CreateTestCase(t *testing.T) {
	mockRepo := new(MockTestCaseRepository)
	service := NewTestCaseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		suiteID := int64(123)
		name := "TestExample"
		classname := "com.example.TestExample"

		mockRepo.On("GetByName", ctx, suiteID, name).Return(nil, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.TestCase")).Return(nil).Once()

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
		assert.Equal(t, domain.ErrInvalidInput, err)
		assert.Nil(t, result)
	})

	t.Run("invalid input - empty classname", func(t *testing.T) {
		result, err := service.CreateTestCase(ctx, 123, "name", "")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidInput, err)
		assert.Nil(t, result)
	})

	t.Run("duplicate test case", func(t *testing.T) {
		suiteID := int64(123)
		name := "TestExample"
		existingTestCase := &domain.TestCase{
			ID:        1,
			SuiteID:   suiteID,
			Name:      name,
			Classname: "com.example.TestExample",
		}

		mockRepo.On("GetByName", ctx, suiteID, name).Return(existingTestCase, nil).Once()

		result, err := service.CreateTestCase(ctx, suiteID, name, "newclassname")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrDuplicateTestCase, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestCaseService_UpdateTestCase(t *testing.T) {
	mockRepo := new(MockTestCaseRepository)
	service := NewTestCaseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		testCaseID := int64(1)
		name := "UpdatedTest"
		classname := "com.example.UpdatedTest"
		existingTestCase := &domain.TestCase{
			ID:        testCaseID,
			SuiteID:   123,
			Name:      "OldTest",
			Classname: "com.example.OldTest",
		}
		updatedTestCase := &domain.TestCase{
			ID:        testCaseID,
			SuiteID:   123,
			Name:      name,
			Classname: classname,
		}

		mockRepo.On("GetByID", ctx, testCaseID).Return(existingTestCase, nil).Once()
		mockRepo.On("Update", ctx, testCaseID, name, classname).Return(updatedTestCase, nil).Once()

		result, err := service.UpdateTestCase(ctx, testCaseID, name, classname)

		assert.NoError(t, err)
		assert.Equal(t, updatedTestCase, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.UpdateTestCase(ctx, 0, "name", "classname")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidInput, err)
		assert.Nil(t, result)
	})

	t.Run("test case not found", func(t *testing.T) {
		testCaseID := int64(999)
		mockRepo.On("GetByID", ctx, testCaseID).Return(nil, nil).Once()

		result, err := service.UpdateTestCase(ctx, testCaseID, "name", "classname")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrTestCaseNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestTestCaseService_DeleteTestCase(t *testing.T) {
	mockRepo := new(MockTestCaseRepository)
	service := NewTestCaseService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		testCaseID := int64(1)
		existingTestCase := &domain.TestCase{
			ID:        testCaseID,
			SuiteID:   123,
			Name:      "TestExample",
			Classname: "com.example.TestExample",
		}

		mockRepo.On("GetByID", ctx, testCaseID).Return(existingTestCase, nil).Once()
		mockRepo.On("Delete", ctx, testCaseID).Return(nil).Once()

		err := service.DeleteTestCase(ctx, testCaseID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		err := service.DeleteTestCase(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidInput, err)
	})

	t.Run("test case not found", func(t *testing.T) {
		testCaseID := int64(999)
		mockRepo.On("GetByID", ctx, testCaseID).Return(nil, nil).Once()

		err := service.DeleteTestCase(ctx, testCaseID)

		assert.Error(t, err)
		assert.Equal(t, domain.ErrTestCaseNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
