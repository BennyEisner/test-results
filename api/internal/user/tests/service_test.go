package application

import (
	"context"
	"testing"
	"time"

	"github.com/BennyEisner/test-results/internal/user/application"
	"github.com/BennyEisner/test-results/internal/user/domain/errors"
	"github.com/BennyEisner/test-results/internal/user/domain/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, id int64, user *models.User) (*models.User, error) {
	args := m.Called(ctx, id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := application.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedUser := &models.User{
			ID:        1,
			Username:  "testuser",
			CreatedAt: time.Now(),
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(expectedUser, nil).Once()

		result, err := service.GetUser(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.GetUser(ctx, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
		assert.Nil(t, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, int64(999)).Return(nil, nil).Once()

		result, err := service.GetUser(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := application.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedUser := &models.User{
			ID:        1,
			Username:  "testuser",
			CreatedAt: time.Now(),
		}

		mockRepo.On("GetByUsername", ctx, "testuser").Return(expectedUser, nil).Once()

		result, err := service.GetUserByUsername(ctx, "testuser")

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.GetUserByUsername(ctx, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username is required")
		assert.Nil(t, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByUsername", ctx, "nonexistent").Return(nil, nil).Once()

		result, err := service.GetUserByUsername(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Equal(t, errors.ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := application.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		username := "newuser"

		mockRepo.On("GetByUsername", ctx, username).Return(nil, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil).Once()

		result, err := service.CreateUser(ctx, username, "test@example.com")

		assert.NoError(t, err)
		assert.Equal(t, username, result.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.CreateUser(ctx, "", "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "username is required")
		assert.Nil(t, result)
	})

	t.Run("duplicate user", func(t *testing.T) {
		username := "existinguser"
		existingUser := &models.User{
			ID:        1,
			Username:  username,
			CreatedAt: time.Now(),
		}

		mockRepo.On("GetByUsername", ctx, username).Return(existingUser, nil).Once()

		result, err := service.CreateUser(ctx, username, "test@example.com")

		assert.Error(t, err)
		assert.Equal(t, errors.ErrUserExists, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := application.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userID := int64(1)
		newUsername := "updateduser"
		existingUser := &models.User{
			ID:        userID,
			Username:  "olduser",
			CreatedAt: time.Now(),
		}
		updatedUser := &models.User{
			ID:        userID,
			Username:  newUsername,
			CreatedAt: existingUser.CreatedAt,
		}

		mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil).Once()
		mockRepo.On("GetByUsername", ctx, newUsername).Return(nil, nil).Once()
		mockRepo.On("Update", ctx, userID, mock.AnythingOfType("*models.User")).Return(updatedUser, nil).Once()

		result, err := service.UpdateUser(ctx, userID, newUsername, "test@example.com")

		assert.NoError(t, err)
		assert.Equal(t, updatedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.UpdateUser(ctx, int64(0), "username", "test@example.com")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
		assert.Nil(t, result)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := int64(999)
		mockRepo.On("GetByID", ctx, userID).Return(nil, nil).Once()

		result, err := service.UpdateUser(ctx, userID, "username", "test@example.com")

		assert.Error(t, err)
		assert.Equal(t, errors.ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := application.NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userID := int64(1)
		existingUser := &models.User{
			ID:        userID,
			Username:  "testuser",
			CreatedAt: time.Now(),
		}

		mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil).Once()
		mockRepo.On("Delete", ctx, userID).Return(nil).Once()

		err := service.DeleteUser(ctx, userID)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		err := service.DeleteUser(ctx, int64(0))

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user ID")
	})

	t.Run("user not found", func(t *testing.T) {
		userID := int64(999)
		mockRepo.On("GetByID", ctx, userID).Return(nil, nil).Once()

		err := service.DeleteUser(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, errors.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
