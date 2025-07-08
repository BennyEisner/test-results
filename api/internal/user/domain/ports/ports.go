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

// UserService defines the interface for user business logic
type UserService interface {
	GetUser(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	CreateUser(ctx context.Context, username, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id int64, username, email string) (*models.User, error)
	DeleteUser(ctx context.Context, id int64) error
}
