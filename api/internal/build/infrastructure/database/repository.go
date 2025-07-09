package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/BennyEisner/test-results/internal/build/domain/models"
	"github.com/BennyEisner/test-results/internal/build/domain/ports"
)

// SQLBuildRepository implements the BuildRepository interface
type SQLBuildRepository struct {
	db *sql.DB
}

// NewSQLBuildRepository creates a new SQL build repository
func NewSQLBuildRepository(db *sql.DB) ports.BuildRepository {
	return &SQLBuildRepository{db: db}
}

// GetByID retrieves a build by its ID
func (r *SQLBuildRepository) GetByID(ctx context.Context, id int64) (*models.Build, error) {
	query := `
		SELECT b.id, b.test_suite_id, ts.project_id, b.build_number, b.ci_provider, b.ci_url, 
		       b.created_at, b.duration, b.test_case_count 
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE b.id = $1
	`

	var build models.Build
	var ciURL sql.NullString
	var duration sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&build.ID, &build.TestSuiteID, &build.ProjectID, &build.BuildNumber,
		&build.CIProvider, &ciURL, &build.CreatedAt,
		&duration, &build.TestCaseCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get build by ID: %w", err)
	}

	if ciURL.Valid {
		build.CIURL = &ciURL.String
	}
	if duration.Valid {
		build.Duration = &duration.Float64
	}

	return &build, nil
}

// GetAllByProjectID retrieves all builds for a project
func (r *SQLBuildRepository) GetAllByProjectID(ctx context.Context, projectID int64) ([]*models.Build, error) {
	query := `
		SELECT b.id, b.test_suite_id, ts.project_id, b.build_number, b.ci_provider, 
		b.ci_url, b.created_at, b.test_case_count, b.duration
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE ts.project_id = $1
		ORDER BY b.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds by project ID: %w", err)
	}
	defer rows.Close()

	var builds []*models.Build
	for rows.Next() {
		var build models.Build
		var ciURL sql.NullString
		var duration sql.NullFloat64

		err := rows.Scan(
			&build.ID, &build.TestSuiteID, &build.ProjectID, &build.BuildNumber,
			&build.CIProvider, &ciURL, &build.CreatedAt,
			&build.TestCaseCount, &duration,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan build: %w", err)
		}

		if ciURL.Valid {
			build.CIURL = &ciURL.String
		}
		if duration.Valid {
			build.Duration = &duration.Float64
		}

		builds = append(builds, &build)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating builds: %w", err)
	}

	return builds, nil
}

// GetAllByTestSuiteID retrieves all builds for a test suite
func (r *SQLBuildRepository) GetAllByTestSuiteID(ctx context.Context, suiteID int64) ([]*models.Build, error) {
	query := `
		SELECT b.id, b.test_suite_id, ts.project_id, b.build_number, b.ci_provider, b.ci_url, 
		       b.created_at, b.duration, b.test_case_count 
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE b.test_suite_id = $1 ORDER BY b.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, suiteID)
	if err != nil {
		return nil, fmt.Errorf("failed to get builds by test suite ID: %w", err)
	}
	defer rows.Close()

	var builds []*models.Build
	for rows.Next() {
		var build models.Build
		var ciURL sql.NullString
		var duration sql.NullFloat64

		err := rows.Scan(
			&build.ID, &build.TestSuiteID, &build.ProjectID, &build.BuildNumber,
			&build.CIProvider, &ciURL, &build.CreatedAt,
			&duration, &build.TestCaseCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan build: %w", err)
		}

		if ciURL.Valid {
			build.CIURL = &ciURL.String
		}
		if duration.Valid {
			build.Duration = &duration.Float64
		}

		builds = append(builds, &build)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating builds: %w", err)
	}

	return builds, nil
}

// Create creates a new build
func (r *SQLBuildRepository) Create(ctx context.Context, build *models.Build) error {
	query := `
		INSERT INTO builds (test_suite_id, build_number, ci_provider, ci_url, 
		                   created_at, test_case_count)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id
	`

	now := time.Now()
	if build.CreatedAt.IsZero() {
		build.CreatedAt = now
	}

	err := r.db.QueryRowContext(ctx, query,
		build.TestSuiteID, build.BuildNumber, build.CIProvider,
		build.CIURL, build.CreatedAt,
		build.TestCaseCount,
	).Scan(&build.ID)
	if err != nil {
		return fmt.Errorf("failed to create build: %w", err)
	}

	return nil
}

// Update updates an existing build
func (r *SQLBuildRepository) Update(ctx context.Context, id int64, build *models.Build) (*models.Build, error) {
	query := `
		UPDATE builds SET test_suite_id = $1, build_number = $2, 
		                 ci_provider = $3, ci_url = $4, test_case_count = $5
		WHERE id = $6
	`

	_, err := r.db.ExecContext(ctx, query,
		build.TestSuiteID, build.BuildNumber, build.CIProvider,
		build.CIURL, build.TestCaseCount, id,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update build: %w", err)
	}

	return r.GetByID(ctx, id)
}

// Delete deletes a build by its ID
func (r *SQLBuildRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM builds WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete build: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("build not found")
	}

	return nil
}
