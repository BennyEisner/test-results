package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/user/domain/errors"
	userconfigmodels "github.com/BennyEisner/test-results/internal/user_config/domain/models"
	userconfigports "github.com/BennyEisner/test-results/internal/user_config/domain/ports"
)

// UserConfigService implements the UserConfigService interface
type UserConfigService struct {
	repo userconfigports.UserConfigRepository
}

func NewUserConfigService(repo userconfigports.UserConfigRepository) userconfigports.UserConfigService {
	return &UserConfigService{repo: repo}
}

func (s *UserConfigService) GetUserConfigs(ctx context.Context, userID int64) ([]*userconfigmodels.UserConfig, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}

	configs, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user configs: %w", err)
	}
	return configs, nil
}

func (s *UserConfigService) GetUserConfig(ctx context.Context, userID int64, key string) (*userconfigmodels.UserConfig, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}
	if key == "" {
		return nil, fmt.Errorf("config key is required")
	}

	config, err := s.repo.GetByUserIDAndKey(ctx, userID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get user config: %w", err)
	}
	if config == nil {
		return nil, errors.ErrConfigNotFound
	}
	return config, nil
}

func (s *UserConfigService) SetUserConfig(ctx context.Context, userID int64, key, value string) (*userconfigmodels.UserConfig, error) {
	if userID <= 0 {
		return nil, fmt.Errorf("invalid user ID")
	}
	if key == "" {
		return nil, fmt.Errorf("config key is required")
	}

	// Check if config already exists
	existingConfig, err := s.repo.GetByUserIDAndKey(ctx, userID, key)
	if err == nil && existingConfig != nil {
		// Update existing config
		existingConfig.Value = value
		updatedConfig, err := s.repo.Update(ctx, existingConfig.ID, existingConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to update user config: %w", err)
		}
		return updatedConfig, nil
	}

	// Create new config
	config := &userconfigmodels.UserConfig{
		UserID: userID,
		Key:    key,
		Value:  value,
	}

	if err := s.repo.Create(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to create user config: %w", err)
	}
	return config, nil
}

func (s *UserConfigService) UpdateUserConfig(ctx context.Context, id int64, value string) (*userconfigmodels.UserConfig, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid config ID")
	}

	config := &userconfigmodels.UserConfig{
		Value: value,
	}

	updatedConfig, err := s.repo.Update(ctx, id, config)
	if err != nil {
		return nil, fmt.Errorf("failed to update user config: %w", err)
	}
	return updatedConfig, nil
}

func (s *UserConfigService) DeleteUserConfig(ctx context.Context, id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid config ID")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete user config: %w", err)
	}
	return nil
}

func (s *UserConfigService) CreateUserConfig(ctx context.Context, userID int64, layouts, activeLayoutID string) (*userconfigmodels.UserConfig, error) {
	config := &userconfigmodels.UserConfig{
		UserID:         userID,
		Layouts:        layouts,
		ActiveLayoutID: activeLayoutID,
	}
	if err := s.repo.Create(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to create user config: %w", err)
	}
	return config, nil
}
