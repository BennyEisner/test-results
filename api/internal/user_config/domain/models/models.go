package models

import "time"

// UserConfig represents user configuration settings
// Supports both key-value and layout-based configs
// (You can remove unused fields if you only need one style)
type UserConfig struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	Key            string    `json:"key,omitempty"`
	Value          string    `json:"value,omitempty"`
	Layouts        string    `json:"layouts,omitempty"` // Storing as a JSON string
	ActiveLayoutID string    `json:"active_layout_id,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
