package application

import (
	"context"
	"fmt"
	"time"

	"github.com/BennyEisner/test-results/internal/domain"
	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// UserService implements the UserService interface
type UserService struct {
	repo ports.UserRepository
}

func NewUserService(repo ports.UserRepository) ports.UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidUsername
	}

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID %d: %w", id, err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	if username == "" {
		return nil, domain.ErrInvalidUsername
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, username string) (*models.User, error) {
	if username == "" {
		return nil, domain.ErrInvalidUsername
	}

	// Check if user already exists
	existingUser, err := s.repo.GetByUsername(ctx, username)
	if err == nil && existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	user := &models.User{
		Username:  username,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, username string) (*models.User, error) {
	if id <= 0 || username == "" {
		return nil, domain.ErrInvalidUsername
	}

	// Check if user exists
	existingUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to check user existence: %w", err)
	}
	if existingUser == nil {
		return nil, domain.ErrUserNotFound
	}

	// Check if new username conflicts with existing user
	conflictingUser, err := s.repo.GetByUsername(ctx, username)
	if err == nil && conflictingUser != nil && conflictingUser.ID != id {
		return nil, domain.ErrUserAlreadyExists
	}

	user := &models.User{
		ID:       id,
		Username: username,
	}

	updatedUser, err := s.repo.Update(ctx, id, user)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return domain.ErrInvalidUsername
	}

	// Check if user exists
	existingUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check user existence: %w", err)
	}
	if existingUser == nil {
		return domain.ErrUserNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
