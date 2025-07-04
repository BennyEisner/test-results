package application

import (
	"context"
	"fmt"
	"time"

	"github.com/BennyEisner/test-results/internal/domain"
)

type UserConfigService struct {
	repo domain.UserConfigRepository
}

func NewUserConfigService(repo domain.UserConfigRepository) domain.UserConfigService {
	return &UserConfigService{repo: repo}
}

func (s *UserConfigService) GetUserConfig(ctx context.Context, userID int) (*domain.UserConfig, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidInput
	}
	config, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user config for user ID %d: %w", userID, err)
	}
	if config == nil {
		return nil, domain.ErrUserConfigNotFound
	}
	return config, nil
}

func (s *UserConfigService) CreateUserConfig(ctx context.Context, userID int, layouts, activeLayoutID string) (*domain.UserConfig, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidInput
	}

	// Check if config already exists
	existingConfig, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check for existing user config: %w", err)
	}
	if existingConfig != nil {
		return nil, fmt.Errorf("user config already exists for user ID %d", userID)
	}

	now := time.Now()
	config := &domain.UserConfig{
		UserID:         userID,
		Layouts:        layouts,
		ActiveLayoutID: activeLayoutID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := s.repo.Create(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to create user config: %w", err)
	}
	return config, nil
}

func (s *UserConfigService) UpdateUserConfig(ctx context.Context, userID int, layouts, activeLayoutID string) (*domain.UserConfig, error) {
	if userID <= 0 {
		return nil, domain.ErrInvalidInput
	}

	// Check if config exists
	existingConfig, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get existing user config: %w", err)
	}
	if existingConfig == nil {
		return nil, domain.ErrUserConfigNotFound
	}

	updatedConfig, err := s.repo.Update(ctx, userID, &domain.UserConfig{
		UserID:         userID,
		Layouts:        layouts,
		ActiveLayoutID: activeLayoutID,
		UpdatedAt:      time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user config: %w", err)
	}
	return updatedConfig, nil
}

func (s *UserConfigService) DeleteUserConfig(ctx context.Context, userID int) error {
	if userID <= 0 {
		return domain.ErrInvalidInput
	}

	// Check if config exists
	existingConfig, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get existing user config: %w", err)
	}
	if existingConfig == nil {
		return domain.ErrUserConfigNotFound
	}

	if err := s.repo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("failed to delete user config: %w", err)
	}
	return nil
}
