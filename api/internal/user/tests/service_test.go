package application

import (
	"context"
	"testing"
	"time"

	"github.com/BennyEisner/test-results/internal/user/domain/models"
	"github.com/BennyEisner/test-results/internal/user/domain/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
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

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, id int, user *models.User) (*models.User, error) {
	args := m.Called(ctx, id, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		expectedUser := &models.User{
			ID:        1,
			Username:  "testuser",
			CreatedAt: time.Now(),
		}

		mockRepo.On("GetByID", ctx, 1).Return(expectedUser, nil).Once()

		result, err := service.GetUserByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.GetUserByID(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidUsername, err)
		assert.Nil(t, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByID", ctx, 999).Return(nil, nil).Once()

		result, err := service.GetUserByID(ctx, 999)

		assert.Error(t, err)
		assert.Equal(t, ports.ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_GetUserByUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
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
		assert.Equal(t, ports.ErrInvalidUsername, err)
		assert.Nil(t, result)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo.On("GetByUsername", ctx, "nonexistent").Return(nil, nil).Once()

		result, err := service.GetUserByUsername(ctx, "nonexistent")

		assert.Error(t, err)
		assert.Equal(t, ports.ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		username := "newuser"

		mockRepo.On("GetByUsername", ctx, username).Return(nil, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil).Once()

		result, err := service.CreateUser(ctx, username)

		assert.NoError(t, err)
		assert.Equal(t, username, result.Username)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.CreateUser(ctx, "")

		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidUsername, err)
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

		result, err := service.CreateUser(ctx, username)

		assert.Error(t, err)
		assert.Equal(t, ports.ErrUserAlreadyExists, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userID := 1
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

		result, err := service.UpdateUser(ctx, userID, newUsername)

		assert.NoError(t, err)
		assert.Equal(t, updatedUser, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid input", func(t *testing.T) {
		result, err := service.UpdateUser(ctx, 0, "username")

		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidUsername, err)
		assert.Nil(t, result)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := 999
		mockRepo.On("GetByID", ctx, userID).Return(nil, nil).Once()

		result, err := service.UpdateUser(ctx, userID, "username")

		assert.Error(t, err)
		assert.Equal(t, ports.ErrUserNotFound, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		userID := 1
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
		err := service.DeleteUser(ctx, 0)

		assert.Error(t, err)
		assert.Equal(t, ports.ErrInvalidUsername, err)
	})

	t.Run("user not found", func(t *testing.T) {
		userID := 999
		mockRepo.On("GetByID", ctx, userID).Return(nil, nil).Once()

		err := service.DeleteUser(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, ports.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}
