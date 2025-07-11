package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/BennyEisner/test-results/internal/auth/domain/models"
	"github.com/BennyEisner/test-results/internal/auth/domain/ports"
)

// SQLAuthRepository implements the AuthRepository interface
type SQLAuthRepository struct {
	db *sql.DB
}

// NewSQLAuthRepository creates a new SQL-based auth repository
func NewSQLAuthRepository(db *sql.DB) ports.AuthRepository {
	return &SQLAuthRepository{db: db}
}

// User operations

func (r *SQLAuthRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO auth_users (provider, provider_id, email, name, first_name, last_name, 
		                       avatar_url, access_token, refresh_token, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		user.Provider, user.ProviderID, user.Email, user.Name, user.FirstName, user.LastName,
		user.AvatarURL, user.AccessToken, user.RefreshToken, user.ExpiresAt, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *SQLAuthRepository) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `
		SELECT id, provider, provider_id, email, name, first_name, last_name, 
		       avatar_url, access_token, refresh_token, expires_at, created_at, updated_at
		FROM auth_users WHERE id = $1`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Provider, &user.ProviderID, &user.Email, &user.Name, &user.FirstName, &user.LastName,
		&user.AvatarURL, &user.AccessToken, &user.RefreshToken, &user.ExpiresAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *SQLAuthRepository) GetUserByProviderID(ctx context.Context, provider, providerID string) (*models.User, error) {
	query := `
		SELECT id, provider, provider_id, email, name, first_name, last_name, 
		       avatar_url, access_token, refresh_token, expires_at, created_at, updated_at
		FROM auth_users WHERE provider = $1 AND provider_id = $2`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, provider, providerID).Scan(
		&user.ID, &user.Provider, &user.ProviderID, &user.Email, &user.Name, &user.FirstName, &user.LastName,
		&user.AvatarURL, &user.AccessToken, &user.RefreshToken, &user.ExpiresAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *SQLAuthRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE auth_users 
		SET email = $1, name = $2, first_name = $3, last_name = $4, 
		    avatar_url = $5, access_token = $6, refresh_token = $7, expires_at = $8, updated_at = $9
		WHERE id = $10`

	result, err := r.db.ExecContext(ctx, query,
		user.Email, user.Name, user.FirstName, user.LastName,
		user.AvatarURL, user.AccessToken, user.RefreshToken, user.ExpiresAt, user.UpdatedAt, user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *SQLAuthRepository) UpsertUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO auth_users (provider, provider_id, email, name, first_name, last_name, 
		                       avatar_url, access_token, refresh_token, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (provider, provider_id) 
		DO UPDATE SET 
			email = EXCLUDED.email,
			name = EXCLUDED.name,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
			avatar_url = EXCLUDED.avatar_url,
			access_token = EXCLUDED.access_token,
			refresh_token = EXCLUDED.refresh_token,
			expires_at = EXCLUDED.expires_at,
			updated_at = EXCLUDED.updated_at
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		user.Provider, user.ProviderID, user.Email, user.Name, user.FirstName, user.LastName,
		user.AvatarURL, user.AccessToken, user.RefreshToken, user.ExpiresAt, user.CreatedAt, user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("failed to upsert user: %w", err)
	}

	return nil
}

// Session operations

func (r *SQLAuthRepository) CreateSession(ctx context.Context, session *models.AuthSession) error {
	query := `
		INSERT INTO auth_sessions (id, user_id, provider, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)`

	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.UserID, session.Provider, session.ExpiresAt, session.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	return nil
}

func (r *SQLAuthRepository) GetSession(ctx context.Context, sessionID string) (*models.AuthSession, error) {
	query := `
		SELECT id, user_id, provider, expires_at, created_at
		FROM auth_sessions WHERE id = $1`

	session := &models.AuthSession{}
	err := r.db.QueryRowContext(ctx, query, sessionID).Scan(
		&session.ID, &session.UserID, &session.Provider, &session.ExpiresAt, &session.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}

func (r *SQLAuthRepository) DeleteSession(ctx context.Context, sessionID string) error {
	query := `DELETE FROM auth_sessions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, sessionID)
	if err != nil {
		return fmt.Errorf("failed to delete session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("session not found")
	}

	return nil
}

func (r *SQLAuthRepository) DeleteExpiredSessions(ctx context.Context) error {
	query := `DELETE FROM auth_sessions WHERE expires_at < $1`

	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to delete expired sessions: %w", err)
	}

	return nil
}

// API Key operations

func (r *SQLAuthRepository) CreateAPIKey(ctx context.Context, apiKey *models.APIKey) error {
	query := `
		INSERT INTO auth_api_keys (user_id, name, key_hash, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	err := r.db.QueryRowContext(ctx, query,
		apiKey.UserID, apiKey.Name, apiKey.KeyHash, apiKey.ExpiresAt, apiKey.CreatedAt, apiKey.UpdatedAt,
	).Scan(&apiKey.ID)

	if err != nil {
		return fmt.Errorf("failed to create API key: %w", err)
	}

	return nil
}

func (r *SQLAuthRepository) GetAPIKeyByHash(ctx context.Context, keyHash string) (*models.APIKey, error) {
	query := `
		SELECT id, user_id, name, key_hash, last_used_at, expires_at, created_at, updated_at
		FROM auth_api_keys WHERE key_hash = $1`

	apiKey := &models.APIKey{}
	err := r.db.QueryRowContext(ctx, query, keyHash).Scan(
		&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.KeyHash, &apiKey.LastUsedAt,
		&apiKey.ExpiresAt, &apiKey.CreatedAt, &apiKey.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("API key not found")
		}
		return nil, fmt.Errorf("failed to get API key: %w", err)
	}

	return apiKey, nil
}

func (r *SQLAuthRepository) UpdateAPIKeyLastUsed(ctx context.Context, keyID int64) error {
	query := `UPDATE auth_api_keys SET last_used_at = $1 WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, time.Now(), keyID)
	if err != nil {
		return fmt.Errorf("failed to update API key last used: %w", err)
	}

	return nil
}

func (r *SQLAuthRepository) DeleteAPIKey(ctx context.Context, keyID int64) error {
	query := `DELETE FROM auth_api_keys WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, keyID)
	if err != nil {
		return fmt.Errorf("failed to delete API key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("API key not found")
	}

	return nil
}

func (r *SQLAuthRepository) ListAPIKeysByUser(ctx context.Context, userID int64) ([]*models.APIKey, error) {
	query := `
		SELECT id, user_id, name, key_hash, last_used_at, expires_at, created_at, updated_at
		FROM auth_api_keys WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query API keys: %w", err)
	}
	defer rows.Close()

	var apiKeys []*models.APIKey
	for rows.Next() {
		apiKey := &models.APIKey{}
		err := rows.Scan(
			&apiKey.ID, &apiKey.UserID, &apiKey.Name, &apiKey.KeyHash, &apiKey.LastUsedAt,
			&apiKey.ExpiresAt, &apiKey.CreatedAt, &apiKey.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan API key: %w", err)
		}
		apiKeys = append(apiKeys, apiKey)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating API keys: %w", err)
	}

	return apiKeys, nil
}
