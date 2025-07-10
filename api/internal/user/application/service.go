package application

import (
	"context"
	"fmt"
	"time"

	"github.com/BennyEisner/test-results/internal/user/domain/errors"
	"github.com/BennyEisner/test-results/internal/user/domain/models"
	"github.com/BennyEisner/test-results/internal/user/domain/ports"
)

// UserService implements the UserService interface
type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*models.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID %d: %w", id, err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	if user == nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, username, email string) (*models.User, error) {
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Check if user already exists
	existingUser, err := s.repo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, errors.ErrUserExists
	}

	user := &models.User{
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int64, username, email string) (*models.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}
	if username == "" {
		return nil, fmt.Errorf("username is required")
	}
	if email == "" {
		return nil, fmt.Errorf("email is required")
	}

	// Check if user exists
	existingUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if existingUser == nil {
		return nil, errors.ErrUserNotFound
	}

	// Check if new username conflicts with existing user
	conflictingUser, err := s.repo.GetByUsername(ctx, username)
	if err == nil && conflictingUser != nil && conflictingUser.ID != id {
		return nil, errors.ErrUserExists
	}

	user := &models.User{
		ID:        id,
		Username:  username,
		Email:     email,
		UpdatedAt: time.Now(),
	}

	updatedUser, err := s.repo.Update(ctx, id, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid user ID")
	}

	// Check if user exists
	existingUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if existingUser == nil {
		return errors.ErrUserNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
