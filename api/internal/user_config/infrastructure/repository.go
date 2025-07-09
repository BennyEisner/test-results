package infrastructure

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

// GetByUserID retrieves a user config by user ID
func (r *SQLUserConfigRepository) GetByUserID(ctx context.Context, userID int64) (*userconfigmodels.UserConfig, error) {
	query := `SELECT id, user_id, layouts, active_layout_id, created_at, updated_at FROM user_configs WHERE user_id = $1`
	var config userconfigmodels.UserConfig
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&config.ID, &config.UserID, &config.Layouts, &config.ActiveLayoutID, &config.CreatedAt, &config.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No config found is not an error
		}
		return nil, fmt.Errorf("failed to get user config: %w", err)
	}
	return &config, nil
}

func (r *SQLUserConfigRepository) GetByUserIDAndKey(ctx context.Context, userID int64, key string) (*userconfigmodels.UserConfig, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *SQLUserConfigRepository) Save(ctx context.Context, config *userconfigmodels.UserConfig) error {
	query := `
        INSERT INTO user_configs (user_id, layouts, active_layout_id, created_at, updated_at)
        VALUES ($1, $2, $3, NOW(), NOW())
        ON CONFLICT (user_id) DO UPDATE SET
            layouts = EXCLUDED.layouts,
            active_layout_id = EXCLUDED.active_layout_id,
            updated_at = NOW()
        RETURNING id, created_at, updated_at
    `
	err := r.db.QueryRowContext(ctx, query, config.UserID, config.Layouts, config.ActiveLayoutID).Scan(&config.ID, &config.CreatedAt, &config.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to save user config: %w", err)
	}
	return nil
}

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
