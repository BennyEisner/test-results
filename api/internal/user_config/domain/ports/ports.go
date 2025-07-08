package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/user_config/domain/models"
)

// UserConfigRepository defines the interface for user config data access
type UserConfigRepository interface {
	GetByUserID(ctx context.Context, userID int64) ([]*models.UserConfig, error)
	GetByUserIDAndKey(ctx context.Context, userID int64, key string) (*models.UserConfig, error)
	Create(ctx context.Context, config *models.UserConfig) error
	Update(ctx context.Context, id int64, config *models.UserConfig) (*models.UserConfig, error)
	Delete(ctx context.Context, id int64) error
}

// UserConfigService defines the interface for user config business logic
type UserConfigService interface {
	GetUserConfigs(ctx context.Context, userID int64) ([]*models.UserConfig, error)
	GetUserConfig(ctx context.Context, userID int64, key string) (*models.UserConfig, error)
	SetUserConfig(ctx context.Context, userID int64, key, value string) (*models.UserConfig, error)
	UpdateUserConfig(ctx context.Context, id int64, value string) (*models.UserConfig, error)
	DeleteUserConfig(ctx context.Context, id int64) error
	CreateUserConfig(ctx context.Context, userID int64, layouts, activeLayoutID string) (*models.UserConfig, error)
}
