package application

import (
	"context"
	"fmt"

	userconfigmodels "github.com/BennyEisner/test-results/internal/user_config/domain/models"
)

// UserConfigService implements the UserConfigService interface
type UserConfigService struct {
	repo interface{}
}

func NewUserConfigService(repo interface{}) interface{} {
	return &UserConfigService{repo: repo}
}

func (s *UserConfigService) GetUserConfigs(ctx context.Context, userID int64) ([]interface{}, error) {
	// Temporarily disabled
	return nil, fmt.Errorf("user config service temporarily disabled")
}

func (s *UserConfigService) GetUserConfig(ctx context.Context, userID int64, key string) (interface{}, error) {
	// Temporarily disabled
	return nil, fmt.Errorf("user config service temporarily disabled")
}

func (s *UserConfigService) SetUserConfig(ctx context.Context, userID int64, key, value string) (interface{}, error) {
	// Temporarily disabled
	return nil, fmt.Errorf("user config service temporarily disabled")
}

func (s *UserConfigService) UpdateUserConfig(ctx context.Context, id int64, value string) (interface{}, error) {
	// Temporarily disabled
	return nil, fmt.Errorf("user config service temporarily disabled")
}

func (s *UserConfigService) DeleteUserConfig(ctx context.Context, id int64) error {
	// Temporarily disabled
	return fmt.Errorf("user config service temporarily disabled")
}
func (s *UserConfigService) CreateUserConfig(ctx context.Context, userID int64, layouts, activeLayoutID string) (*userconfigmodels.UserConfig, error) {
	// Dummy implementation for now
	return &userconfigmodels.UserConfig{
		UserID:         userID,
		Layouts:        layouts,
		ActiveLayoutID: activeLayoutID,
	}, nil
}
