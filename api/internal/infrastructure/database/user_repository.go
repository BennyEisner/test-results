package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// SQLUserRepository implements the UserRepository interface
type SQLUserRepository struct {
	db *sql.DB
}

// NewSQLUserRepository creates a new SQL user repository
func NewSQLUserRepository(db *sql.DB) ports.UserRepository {
	return &SQLUserRepository{db: db}
}

// GetByID retrieves a user by its ID
func (r *SQLUserRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, username, created_at FROM users WHERE id = $1`

	var user models.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return &user, nil
}

// GetByUsername retrieves a user by its username
func (r *SQLUserRepository) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT id, username, created_at FROM users WHERE username = $1`

	var user models.User

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

// Create creates a new user
func (r *SQLUserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, created_at) VALUES ($1, $2) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		user.Username, user.CreatedAt,
	).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update updates an existing user
func (r *SQLUserRepository) Update(ctx context.Context, id int, user *models.User) (*models.User, error) {
	query := `UPDATE users SET username = $1 WHERE id = $2 RETURNING id, username, created_at`

	var updatedUser models.User

	err := r.db.QueryRowContext(ctx, query, user.Username, id).Scan(
		&updatedUser.ID, &updatedUser.Username, &updatedUser.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &updatedUser, nil
}

// Delete deletes a user by its ID
func (r *SQLUserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}
