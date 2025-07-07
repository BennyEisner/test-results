package models

import "time"

// User represents a user in the system
type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserConfig represents user configuration
type UserConfig struct {
	ID     int64  `json:"id"`
	UserID int64  `json:"user_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
