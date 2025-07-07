package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/user/domain/models"
	"github.com/BennyEisner/test-results/internal/user/domain/ports"
)

// SQLUserConfigRepository implements the UserConfigRepository interface
type SQLUserConfigRepository struct {
	db *sql.DB
}

// NewSQLUserConfigRepository creates a new SQL user config repository
func NewSQLUserConfigRepository(db *sql.DB) ports.UserConfigRepository {
	return &SQLUserConfigRepository{db: db}
}

// GetByUserID retrieves all configs for a user
func (r *SQLUserConfigRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.UserConfig, error) {
	query := `SELECT id, user_id, key, value FROM user_configs WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user configs: %w", err)
	}
	defer rows.Close()

	var configs []*models.UserConfig
	for rows.Next() {
		var config models.UserConfig
		err := rows.Scan(&config.ID, &config.UserID, &config.Key, &config.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user config: %w", err)
		}
		configs = append(configs, &config)
	}

	return configs, nil
}

// GetByUserIDAndKey retrieves a specific config for a user
func (r *SQLUserConfigRepository) GetByUserIDAndKey(ctx context.Context, userID int64, key string) (*models.UserConfig, error) {
	query := `SELECT id, user_id, key, value FROM user_configs WHERE user_id = $1 AND key = $2`

	var config models.UserConfig

	err := r.db.QueryRowContext(ctx, query, userID, key).Scan(
		&config.ID, &config.UserID, &config.Key, &config.Value,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user config: %w", err)
	}

	return &config, nil
}

// Create creates a new user config
func (r *SQLUserConfigRepository) Create(ctx context.Context, config *models.UserConfig) error {
	query := `INSERT INTO user_configs (user_id, key, value) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		config.UserID, config.Key, config.Value,
	).Scan(&config.ID)
	if err != nil {
		return fmt.Errorf("failed to create user config: %w", err)
	}

	return nil
}

// Update updates an existing user config
func (r *SQLUserConfigRepository) Update(ctx context.Context, id int64, config *models.UserConfig) (*models.UserConfig, error) {
	query := `UPDATE user_configs SET value = $1 WHERE id = $2 RETURNING id, user_id, key, value`

	var updatedConfig models.UserConfig

	err := r.db.QueryRowContext(ctx, query, config.Value, id).Scan(
		&updatedConfig.ID, &updatedConfig.UserID, &updatedConfig.Key, &updatedConfig.Value,
	)
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
