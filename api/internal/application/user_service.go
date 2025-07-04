package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
)

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	if id <= 0 {
		return nil, domain.ErrInvalidInput
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

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	if username == "" {
		return nil, domain.ErrInvalidInput
	}
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by username %s: %w", username, err)
	}
	if user == nil {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, username string) (*domain.User, error) {
	if username == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if user already exists
	existingUser, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing user: %w", err)
	}
	if existingUser != nil {
		return nil, domain.ErrDuplicateUser
	}

	user := &domain.User{
		Username: username,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int, username string) (*domain.User, error) {
	if id <= 0 || username == "" {
		return nil, domain.ErrInvalidInput
	}

	// Check if user exists
	existingUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user: %w", err)
	}
	if existingUser == nil {
		return nil, domain.ErrUserNotFound
	}

	// Check if new username already exists (if different from current)
	if existingUser.Username != username {
		userWithNewUsername, err := s.repo.GetByUsername(ctx, username)
		if err != nil {
			return nil, fmt.Errorf("failed to check for username conflict: %w", err)
		}
		if userWithNewUsername != nil {
			return nil, domain.ErrDuplicateUser
		}
	}

	updatedUser, err := s.repo.Update(ctx, id, &domain.User{Username: username})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	return updatedUser, nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	if id <= 0 {
		return domain.ErrInvalidInput
	}

	// Check if user exists
	existingUser, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get existing user: %w", err)
	}
	if existingUser == nil {
		return domain.ErrUserNotFound
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
