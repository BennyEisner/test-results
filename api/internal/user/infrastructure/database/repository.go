package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/user/domain/models"
	"github.com/BennyEisner/test-results/internal/user/domain/ports"
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
func (r *SQLUserRepository) GetByID(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE id = $1`

	var user models.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt,
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
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE username = $1`

	var user models.User

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

// GetByEmail retrieves a user by its email
func (r *SQLUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, email, created_at, updated_at FROM users WHERE email = $1`

	var user models.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return &user, nil
}

// Create creates a new user
func (r *SQLUserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (username, email, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		user.Username, user.Email, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// Update updates an existing user
func (r *SQLUserRepository) Update(ctx context.Context, id int64, user *models.User) (*models.User, error) {
	query := `UPDATE users SET username = $1, email = $2, updated_at = $3 WHERE id = $4 RETURNING id, username, email, created_at, updated_at`

	var updatedUser models.User

	err := r.db.QueryRowContext(ctx, query, user.Username, user.Email, user.UpdatedAt, id).Scan(
		&updatedUser.ID, &updatedUser.Username, &updatedUser.Email, &updatedUser.CreatedAt, &updatedUser.UpdatedAt,
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
func (r *SQLUserRepository) Delete(ctx context.Context, id int64) error {
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
