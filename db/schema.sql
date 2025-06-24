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
    test_case_count INTEGER
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
