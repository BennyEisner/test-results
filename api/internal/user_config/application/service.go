package application

import (
	"context"
	"fmt"

	userconfigmodels "github.com/BennyEisner/test-results/internal/user_config/domain/models"
	"github.com/BennyEisner/test-results/internal/user_config/domain/ports"
)

// UserConfigService implements the UserConfigService interface
type UserConfigService struct {
	repo ports.UserConfigRepository
}

func NewUserConfigService(repo ports.UserConfigRepository) ports.UserConfigService {
	return &UserConfigService{repo: repo}
}

func (s *UserConfigService) GetUserConfigs(ctx context.Context, userID int64) ([]*userconfigmodels.UserConfig, error) {
	config, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if config == nil {
		return []*userconfigmodels.UserConfig{}, nil
	}
	return []*userconfigmodels.UserConfig{config}, nil
}

func (s *UserConfigService) GetUserConfig(ctx context.Context, userID int64, key string) (*userconfigmodels.UserConfig, error) {
	return nil, fmt.Errorf("not implemented")
}

func (s *UserConfigService) SaveUserConfig(ctx context.Context, userID int64, layouts, activeLayoutID string) (*userconfigmodels.UserConfig, error) {
	config, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		// Assuming an error here means something went wrong with the query, not that the config doesn't exist.
		return nil, fmt.Errorf("failed to get user config: %w", err)
	}

	if config == nil {
		// Config does not exist, create a new one
		config = &userconfigmodels.UserConfig{
			UserID:         userID,
			Layouts:        layouts,
			ActiveLayoutID: activeLayoutID,
		}
	} else {
		// Config exists, update it
		config.Layouts = layouts
		config.ActiveLayoutID = activeLayoutID
	}

	if err := s.repo.Save(ctx, config); err != nil {
		return nil, fmt.Errorf("failed to save user config: %w", err)
	}
	return config, nil
}

func (s *UserConfigService) DeleteUserConfig(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *UserConfigService) UpdateActiveLayoutID(ctx context.Context, userID int64, activeLayoutID string) error {
	config, err := s.repo.GetByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user config: %w", err)
	}
	if config == nil {
		return fmt.Errorf("user config not found for user %d", userID)
	}

	config.ActiveLayoutID = activeLayoutID
	if err := s.repo.Save(ctx, config); err != nil {
		return fmt.Errorf("failed to update active layout ID: %w", err)
	}
	return nil
}
