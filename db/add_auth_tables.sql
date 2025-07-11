-- Authentication tables for OAuth2 and API key authentication

-- Users table for authenticated users
CREATE TABLE auth_users (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(50) NOT NULL,           -- "okta", "github", etc.
    provider_id VARCHAR(255) NOT NULL,       -- External provider's user ID
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    avatar_url TEXT,
    access_token TEXT,                       -- Encrypted OAuth2 access token
    refresh_token TEXT,                      -- Encrypted OAuth2 refresh token
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_id),
    UNIQUE(email)
);

-- API Keys table for CLI/Jenkins authentication
CREATE TABLE auth_api_keys (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,              -- Human-readable name for the key
    key_hash VARCHAR(255) NOT NULL UNIQUE,   -- Hashed API key
    last_used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Sessions table for web authentication
CREATE TABLE auth_sessions (
    id VARCHAR(255) PRIMARY KEY,             -- Session ID (UUID)
    user_id INTEGER NOT NULL REFERENCES auth_users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_auth_users_provider_id ON auth_users(provider, provider_id);
CREATE INDEX idx_auth_users_email ON auth_users(email);
CREATE INDEX idx_auth_api_keys_hash ON auth_api_keys(key_hash);
CREATE INDEX idx_auth_api_keys_user_id ON auth_api_keys(user_id);
CREATE INDEX idx_auth_sessions_user_id ON auth_sessions(user_id);
CREATE INDEX idx_auth_sessions_expires_at ON auth_sessions(expires_at);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Triggers to automatically update updated_at
CREATE TRIGGER update_auth_users_updated_at BEFORE UPDATE ON auth_users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_auth_api_keys_updated_at BEFORE UPDATE ON auth_api_keys
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column(); 