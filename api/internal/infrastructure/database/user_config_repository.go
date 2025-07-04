package database

import (
	"context"
	"database/sql"

	"github.com/BennyEisner/test-results/internal/domain"
)

type SQLUserConfigRepository struct {
	db *sql.DB
}

func NewSQLUserConfigRepository(db *sql.DB) domain.UserConfigRepository {
	return &SQLUserConfigRepository{db: db}
}

func (r *SQLUserConfigRepository) GetByUserID(ctx context.Context, userID int) (*domain.UserConfig, error) {
	query := `SELECT id, user_id, layouts, active_layout_id, created_at, updated_at FROM user_configs WHERE user_id = $1`
	config := &domain.UserConfig{}
	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&config.ID, &config.UserID, &config.Layouts, &config.ActiveLayoutID, &config.CreatedAt, &config.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return config, nil
}

func (r *SQLUserConfigRepository) Create(ctx context.Context, config *domain.UserConfig) error {
	query := `INSERT INTO user_configs (user_id, layouts, active_layout_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.db.QueryRowContext(ctx, query, config.UserID, config.Layouts, config.ActiveLayoutID, config.CreatedAt, config.UpdatedAt).Scan(&config.ID)
}

func (r *SQLUserConfigRepository) Update(ctx context.Context, userID int, config *domain.UserConfig) (*domain.UserConfig, error) {
	query := `UPDATE user_configs SET layouts = $1, active_layout_id = $2, updated_at = $3 WHERE user_id = $4 RETURNING id, user_id, layouts, active_layout_id, created_at, updated_at`
	updatedConfig := &domain.UserConfig{}
	if err := r.db.QueryRowContext(ctx, query, config.Layouts, config.ActiveLayoutID, config.UpdatedAt, userID).Scan(&updatedConfig.ID, &updatedConfig.UserID, &updatedConfig.Layouts, &updatedConfig.ActiveLayoutID, &updatedConfig.CreatedAt, &updatedConfig.UpdatedAt); err != nil {
		return nil, err
	}
	return updatedConfig, nil
}

func (r *SQLUserConfigRepository) Delete(ctx context.Context, userID int) error {
	query := `DELETE FROM user_configs WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
