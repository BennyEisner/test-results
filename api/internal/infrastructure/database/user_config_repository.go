package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
)

// SQLUserConfigRepository implements the UserConfigRepository interface
type SQLUserConfigRepository struct {
	db *sql.DB
}

// NewSQLUserConfigRepository creates a new SQL user config repository
func NewSQLUserConfigRepository(db *sql.DB) ports.UserConfigRepository {
	return &SQLUserConfigRepository{db: db}
}

// GetByUserID retrieves a user config by user ID
func (r *SQLUserConfigRepository) GetByUserID(ctx context.Context, userID int) (*models.UserConfig, error) {
	query := `SELECT id, user_id, layouts, active_layout_id, created_at, updated_at FROM user_configs WHERE user_id = $1`

	var config models.UserConfig

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&config.ID, &config.UserID, &config.Layouts, &config.ActiveLayoutID, &config.CreatedAt, &config.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user config by user ID: %w", err)
	}

	return &config, nil
}

// Create creates a new user config
func (r *SQLUserConfigRepository) Create(ctx context.Context, config *models.UserConfig) error {
	query := `INSERT INTO user_configs (user_id, layouts, active_layout_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		config.UserID, config.Layouts, config.ActiveLayoutID, config.CreatedAt, config.UpdatedAt,
	).Scan(&config.ID)
	if err != nil {
		return fmt.Errorf("failed to create user config: %w", err)
	}

	return nil
}

// Update updates an existing user config
func (r *SQLUserConfigRepository) Update(ctx context.Context, userID int, config *models.UserConfig) (*models.UserConfig, error) {
	query := `UPDATE user_configs SET layouts = $1, active_layout_id = $2, updated_at = $3 WHERE user_id = $4 RETURNING id, user_id, layouts, active_layout_id, created_at, updated_at`

	var updatedConfig models.UserConfig

	err := r.db.QueryRowContext(ctx, query, config.Layouts, config.ActiveLayoutID, config.UpdatedAt, userID).Scan(
		&updatedConfig.ID, &updatedConfig.UserID, &updatedConfig.Layouts, &updatedConfig.ActiveLayoutID, &updatedConfig.CreatedAt, &updatedConfig.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update user config: %w", err)
	}

	return &updatedConfig, nil
}

// Delete deletes a user config by user ID
func (r *SQLUserConfigRepository) Delete(ctx context.Context, userID int) error {
	query := `DELETE FROM user_configs WHERE user_id = $1`

	result, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user config: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user config not found")
	}

	return nil
}
