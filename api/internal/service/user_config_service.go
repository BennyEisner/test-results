package service

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/BennyEisner/test-results/internal/models"
)

type UserConfigService struct {
	DB *sql.DB
}

func NewUserConfigService(db *sql.DB) *UserConfigService {
	return &UserConfigService{DB: db}
}

// Marshal layouts interface to JSON to be properly read by the database
func (s *UserConfigService) GetUserConfig(userID int) (*models.UserConfig, error) {
	row := s.DB.QueryRow("SELECT id, user_id, layouts, active_layout_id, created_at, updated_at FROM user_configs WHERE user_id = $1", userID)

	config := &models.UserConfig{}
	err := row.Scan(&config.ID, &config.UserID, &config.Layouts, &config.ActiveLayoutID, &config.CreatedAt, &config.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No config found is not an error
		}
		return nil, err
	}
	return config, nil
}

func (s *UserConfigService) SaveUserConfig(userID int, layouts interface{}, activeLayoutID string) error {
	layoutsJSON, err := json.Marshal(layouts)
	if err != nil {
		return err
	}

	query := `
        INSERT INTO user_configs (user_id, layouts, active_layout_id, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id) DO UPDATE SET
            layouts = EXCLUDED.layouts,
            active_layout_id = EXCLUDED.active_layout_id,
            updated_at = EXCLUDED.updated_at
    `
	_, err = s.DB.Exec(query, userID, layoutsJSON, activeLayoutID, time.Now(), time.Now())
	return err
}
