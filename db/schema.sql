-- Table: projects
CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

-- Table: builds
CREATE TABLE builds (
    id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    build_number TEXT NOT NULL,
    ci_provider TEXT NOT NULL,
    ci_url TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

-- Table: test_suites
CREATE TABLE test_suites (
    id SERIAL PRIMARY KEY,
    build_id INTEGER NOT NULL REFERENCES builds(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    parent_id INTEGER REFERENCES test_suites(id) ON DELETE CASCADE,
    time DOUBLE PRECISION NOT NULL
);

-- Table: test_cases
CREATE TABLE test_cases (
    id SERIAL PRIMARY KEY,
    suite_id INTEGER NOT NULL REFERENCES test_suites(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    classname TEXT NOT NULL,
    time DOUBLE PRECISION NOT NULL,
    status TEXT NOT NULL DEFAULT 'passed'  -- ENUM-like: passed, failed, skipped
);

-- Table: failures
CREATE TABLE failures (
    id SERIAL PRIMARY KEY,
    test_case_id INTEGER NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
    message TEXT,
    type TEXT,
    details TEXT
);

-- Indexes for performance (optional but recommended)
CREATE INDEX idx_builds_project_id ON builds(project_id);
CREATE INDEX idx_test_suites_build_id ON test_suites(build_id);
CREATE INDEX idx_test_suites_parent_id ON test_suites(parent_id);
CREATE INDEX idx_test_cases_suite_id ON test_cases(suite_id);
CREATE INDEX idx_failures_test_case_id ON failures(test_case_id);
