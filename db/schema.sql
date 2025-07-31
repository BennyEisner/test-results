-- Table: projects
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Table: test_suites
CREATE TABLE test_suites (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    parent_id INTEGER REFERENCES test_suites(id) ON DELETE CASCADE,
    time DOUBLE PRECISION NOT NULL
);

-- Table: builds
CREATE TABLE builds (
    id SERIAL PRIMARY KEY,
    test_suite_id INTEGER NOT NULL REFERENCES test_suites(id) ON DELETE CASCADE,
    build_number TEXT NOT NULL,
    ci_provider TEXT NOT NULL,
    ci_url TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    test_case_count INTEGER,
    duration DOUBLE PRECISION
);

-- Table: test_cases
CREATE TABLE test_cases (
    id SERIAL PRIMARY KEY,
    suite_id INTEGER NOT NULL REFERENCES test_suites(id) ON DELETE CASCADE, -- Defines which suite this test case belongs to
    name TEXT NOT NULL,
    classname TEXT NOT NULL
);

-- Table: build_test_case_executions
CREATE TABLE build_test_case_executions (
    id SERIAL PRIMARY KEY,
    build_id INTEGER NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
    test_case_id INTEGER NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    status TEXT NOT NULL, -- e.g., 'passed', 'failed', 'skipped', 'error'
    execution_time DOUBLE PRECISION, -- Actual time taken for this specific execution
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (build_id, test_case_id) -- Ensures one record per test case per build
);

-- Table: failures
CREATE TABLE failures (
    id SERIAL PRIMARY KEY,
    build_test_case_execution_id INTEGER NOT NULL REFERENCES build_test_case_executions(id) ON DELETE CASCADE,
    message TEXT,
    type TEXT,
    details TEXT,
    UNIQUE (build_test_case_execution_id) -- Assuming one failure detail entry per execution
);

-- Indexes for performance (optional but recommended)
CREATE INDEX idx_test_suites_project_id ON test_suites(project_id);
CREATE INDEX idx_builds_test_suite_id ON builds(test_suite_id);
CREATE INDEX idx_test_suites_parent_id ON test_suites(parent_id);
CREATE INDEX idx_test_cases_suite_id ON test_cases(suite_id);
CREATE INDEX idx_btexec_build_id ON build_test_case_executions(build_id);
CREATE INDEX idx_btexec_test_case_id ON build_test_case_executions(test_case_id);
CREATE INDEX idx_failures_btexec_id ON failures(build_test_case_execution_id);
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