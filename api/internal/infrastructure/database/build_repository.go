package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/BennyEisner/test-results/internal/domain/models"
	"github.com/BennyEisner/test-results/internal/domain/ports"
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
		SELECT id, test_suite_id, project_id, build_number, ci_provider, ci_url, 
		       created_at, started_at, ended_at, duration, test_case_count 
		FROM builds WHERE id = $1
	`

	var build models.Build
	var ciURL sql.NullString
	var startedAt, endedAt sql.NullTime
	var duration sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&build.ID, &build.TestSuiteID, &build.ProjectID, &build.BuildNumber,
		&build.CIProvider, &ciURL, &build.CreatedAt, &startedAt, &endedAt,
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
	if startedAt.Valid {
		build.StartedAt = &startedAt.Time
	}
	if endedAt.Valid {
		build.EndedAt = &endedAt.Time
	}
	if duration.Valid {
		build.Duration = &duration.Float64
	}

	return &build, nil
}

// GetAllByProjectID retrieves all builds for a project
func (r *SQLBuildRepository) GetAllByProjectID(ctx context.Context, projectID int64) ([]*models.Build, error) {
	query := `
		SELECT id, test_suite_id, project_id, build_number, ci_provider, ci_url, 
		       created_at, started_at, ended_at, duration, test_case_count 
		FROM builds WHERE project_id = $1 ORDER BY created_at DESC
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
		var startedAt, endedAt sql.NullTime
		var duration sql.NullFloat64

		err := rows.Scan(
			&build.ID, &build.TestSuiteID, &build.ProjectID, &build.BuildNumber,
			&build.CIProvider, &ciURL, &build.CreatedAt, &startedAt, &endedAt,
			&duration, &build.TestCaseCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan build: %w", err)
		}

		if ciURL.Valid {
			build.CIURL = &ciURL.String
		}
		if startedAt.Valid {
			build.StartedAt = &startedAt.Time
		}
		if endedAt.Valid {
			build.EndedAt = &endedAt.Time
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
		SELECT id, test_suite_id, project_id, build_number, ci_provider, ci_url, 
		       created_at, started_at, ended_at, duration, test_case_count 
		FROM builds WHERE test_suite_id = $1 ORDER BY created_at DESC
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
		var startedAt, endedAt sql.NullTime
		var duration sql.NullFloat64

		err := rows.Scan(
			&build.ID, &build.TestSuiteID, &build.ProjectID, &build.BuildNumber,
			&build.CIProvider, &ciURL, &build.CreatedAt, &startedAt, &endedAt,
			&duration, &build.TestCaseCount,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan build: %w", err)
		}

		if ciURL.Valid {
			build.CIURL = &ciURL.String
		}
		if startedAt.Valid {
			build.StartedAt = &startedAt.Time
		}
		if endedAt.Valid {
			build.EndedAt = &endedAt.Time
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
		INSERT INTO builds (test_suite_id, project_id, build_number, ci_provider, ci_url, 
		                   created_at, started_at, ended_at, duration, test_case_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id
	`

	now := time.Now()
	if build.CreatedAt.IsZero() {
		build.CreatedAt = now
	}

	err := r.db.QueryRowContext(ctx, query,
		build.TestSuiteID, build.ProjectID, build.BuildNumber, build.CIProvider,
		build.CIURL, build.CreatedAt, build.StartedAt, build.EndedAt,
		build.Duration, build.TestCaseCount,
	).Scan(&build.ID)
	if err != nil {
		return fmt.Errorf("failed to create build: %w", err)
	}

	return nil
}

// Update updates an existing build
func (r *SQLBuildRepository) Update(ctx context.Context, id int64, build *models.Build) (*models.Build, error) {
	query := `
		UPDATE builds SET test_suite_id = $1, project_id = $2, build_number = $3, 
		                 ci_provider = $4, ci_url = $5, started_at = $6, ended_at = $7, 
		                 duration = $8, test_case_count = $9
		WHERE id = $10 RETURNING id, test_suite_id, project_id, build_number, ci_provider, 
		                        ci_url, created_at, started_at, ended_at, duration, test_case_count
	`

	var updatedBuild models.Build
	var ciURL sql.NullString
	var startedAt, endedAt sql.NullTime
	var duration sql.NullFloat64

	err := r.db.QueryRowContext(ctx, query,
		build.TestSuiteID, build.ProjectID, build.BuildNumber, build.CIProvider,
		build.CIURL, build.StartedAt, build.EndedAt, build.Duration,
		build.TestCaseCount, id,
	).Scan(
		&updatedBuild.ID, &updatedBuild.TestSuiteID, &updatedBuild.ProjectID,
		&updatedBuild.BuildNumber, &updatedBuild.CIProvider, &ciURL,
		&updatedBuild.CreatedAt, &startedAt, &endedAt, &duration, &updatedBuild.TestCaseCount,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to update build: %w", err)
	}

	if ciURL.Valid {
		updatedBuild.CIURL = &ciURL.String
	}
	if startedAt.Valid {
		updatedBuild.StartedAt = &startedAt.Time
	}
	if endedAt.Valid {
		updatedBuild.EndedAt = &endedAt.Time
	}
	if duration.Valid {
		updatedBuild.Duration = &duration.Float64
	}

	return &updatedBuild, nil
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
