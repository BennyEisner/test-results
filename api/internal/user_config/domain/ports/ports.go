package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/user_config/domain/models"
)

// UserConfigRepository defines the interface for user config data access
type UserConfigRepository interface {
	GetByUserID(ctx context.Context, userID int64) (*models.UserConfig, error)
	GetByUserIDAndKey(ctx context.Context, userID int64, key string) (*models.UserConfig, error)
	Save(ctx context.Context, config *models.UserConfig) error
	Delete(ctx context.Context, id int64) error
}

// UserConfigService defines the interface for user config business logic
type UserConfigService interface {
	GetUserConfigs(ctx context.Context, userID int64) ([]*models.UserConfig, error)
	GetUserConfig(ctx context.Context, userID int64, key string) (*models.UserConfig, error)
	SaveUserConfig(ctx context.Context, userID int64, layouts, activeLayoutID string) (*models.UserConfig, error)
	DeleteUserConfig(ctx context.Context, id int64) error
	UpdateActiveLayoutID(ctx context.Context, userID int64, activeLayoutID string) error
}
