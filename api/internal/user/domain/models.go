package models

import "time"

// User represents a user in the system
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

// UserConfig represents user configuration settings
type UserConfig struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	Layouts        string    `json:"layouts"` // Storing as a JSON string
	ActiveLayoutID string    `json:"active_layout_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
