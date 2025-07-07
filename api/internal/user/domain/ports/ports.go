package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/user/domain/models"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByUsername(ctx context.Context, username string) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Create(ctx context.Context, user *models.User) error
	Update(ctx context.Context, id int64, user *models.User) (*models.User, error)
	Delete(ctx context.Context, id int64) error
}

// UserConfigRepository defines the interface for user config data access
type UserConfigRepository interface {
	GetByUserID(ctx context.Context, userID int64) ([]*models.UserConfig, error)
	GetByUserIDAndKey(ctx context.Context, userID int64, key string) (*models.UserConfig, error)
	Create(ctx context.Context, config *models.UserConfig) error
	Update(ctx context.Context, id int64, config *models.UserConfig) (*models.UserConfig, error)
	Delete(ctx context.Context, id int64) error
}

// UserService defines the interface for user business logic
type UserService interface {
	GetUser(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, username, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, username, email string) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
}

// UserConfigService defines the interface for user config business logic
type UserConfigService interface {
	GetUserConfigs(ctx context.Context, userID int64) ([]*models.UserConfig, error)
	GetUserConfig(ctx context.Context, userID int64, key string) (*models.UserConfig, error)
	SetUserConfig(ctx context.Context, userID int64, key, value string) (*models.UserConfig, error)
	UpdateUserConfig(ctx context.Context, id int64, value string) (*models.UserConfig, error)
	DeleteUserConfig(ctx context.Context, id int64) error
}
