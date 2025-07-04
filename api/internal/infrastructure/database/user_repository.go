package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/BennyEisner/test-results/internal/domain"
)

type SQLUserRepository struct {
	db *sql.DB
}

func NewSQLUserRepository(db *sql.DB) domain.UserRepository {
	return &SQLUserRepository{db: db}
}

func (r *SQLUserRepository) GetByID(ctx context.Context, id int) (*domain.User, error) {
	query := `SELECT id, username, created_at FROM users WHERE id = $1`
	user := &domain.User{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Username, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *SQLUserRepository) GetByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `SELECT id, username, created_at FROM users WHERE username = $1`
	user := &domain.User{}
	if err := r.db.QueryRowContext(ctx, query, username).Scan(&user.ID, &user.Username, &user.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *SQLUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, created_at) VALUES ($1, $2) RETURNING id`
	now := time.Now()
	user.CreatedAt = now
	return r.db.QueryRowContext(ctx, query, user.Username, now).Scan(&user.ID)
}

func (r *SQLUserRepository) Update(ctx context.Context, id int, user *domain.User) (*domain.User, error) {
	query := `UPDATE users SET username = $1 WHERE id = $2 RETURNING id, username, created_at`
	updatedUser := &domain.User{}
	if err := r.db.QueryRowContext(ctx, query, user.Username, id).Scan(&updatedUser.ID, &updatedUser.Username, &updatedUser.CreatedAt); err != nil {
		return nil, err
	}
	return updatedUser, nil
}

func (r *SQLUserRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
