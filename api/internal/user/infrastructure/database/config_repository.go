package database

import (
	"context"
	"database/sql"
	"fmt"

	userconfigmodels "github.com/BennyEisner/test-results/internal/user_config/domain/models"
	userconfigports "github.com/BennyEisner/test-results/internal/user_config/domain/ports"
)

// SQLUserConfigRepository implements the UserConfigRepository interface
type SQLUserConfigRepository struct {
	db *sql.DB
}

// NewSQLUserConfigRepository creates a new SQL user config repository
func NewSQLUserConfigRepository(db *sql.DB) userconfigports.UserConfigRepository {
	return &SQLUserConfigRepository{db: db}
}

// GetByUserID retrieves all configs for a user
func (r *SQLUserConfigRepository) GetByUserID(ctx context.Context, userID int64) ([]*userconfigmodels.UserConfig, error) {
	query := `SELECT id, user_id, key, value, layouts, active_layout_id, created_at, updated_at FROM user_configs WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user configs: %w", err)
	}
	defer rows.Close()
	var configs []*userconfigmodels.UserConfig
	for rows.Next() {
		var config userconfigmodels.UserConfig
		err := rows.Scan(&config.ID, &config.UserID, &config.Key, &config.Value, &config.Layouts, &config.ActiveLayoutID, &config.CreatedAt, &config.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user config: %w", err)
		}
		configs = append(configs, &config)
	}
	return configs, nil
}

// GetByUserIDAndKey retrieves a specific config for a user
func (r *SQLUserConfigRepository) GetByUserIDAndKey(ctx context.Context, userID int64, key string) (*userconfigmodels.UserConfig, error) {
	query := `SELECT id, user_id, key, value, layouts, active_layout_id, created_at, updated_at FROM user_configs WHERE user_id = $1 AND key = $2`
	var config userconfigmodels.UserConfig
	err := r.db.QueryRowContext(ctx, query, userID, key).Scan(&config.ID, &config.UserID, &config.Key, &config.Value, &config.Layouts, &config.ActiveLayoutID, &config.CreatedAt, &config.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user config: %w", err)
	}
	return &config, nil
}

// Create creates a new user config
func (r *SQLUserConfigRepository) Create(ctx context.Context, config *userconfigmodels.UserConfig) error {
	query := `INSERT INTO user_configs (user_id, key, value, layouts, active_layout_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRowContext(ctx, query, config.UserID, config.Key, config.Value, config.Layouts, config.ActiveLayoutID, config.CreatedAt, config.UpdatedAt).Scan(&config.ID)
	if err != nil {
		return fmt.Errorf("failed to create user config: %w", err)
	}
	return nil
}

// Update updates an existing user config
func (r *SQLUserConfigRepository) Update(ctx context.Context, id int64, config *userconfigmodels.UserConfig) (*userconfigmodels.UserConfig, error) {
	query := `UPDATE user_configs SET key = $1, value = $2, layouts = $3, active_layout_id = $4, updated_at = $5 WHERE id = $6 RETURNING id, user_id, key, value, layouts, active_layout_id, created_at, updated_at`
	var updatedConfig userconfigmodels.UserConfig
	err := r.db.QueryRowContext(ctx, query, config.Key, config.Value, config.Layouts, config.ActiveLayoutID, config.UpdatedAt, id).Scan(&updatedConfig.ID, &updatedConfig.UserID, &updatedConfig.Key, &updatedConfig.Value, &updatedConfig.Layouts, &updatedConfig.ActiveLayoutID, &updatedConfig.CreatedAt, &updatedConfig.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update user config: %w", err)
	}
	return &updatedConfig, nil
}

// Delete deletes a user config by its ID
func (r *SQLUserConfigRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM user_configs WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
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
