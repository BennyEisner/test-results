package database

import (
	"context"
	"database/sql"

	"github.com/BennyEisner/test-results/internal/domain"
)

type SQLBuildRepository struct {
	db *sql.DB
}

func NewSQLBuildRepository(db *sql.DB) domain.BuildRepository {
	return &SQLBuildRepository{db: db}
}

func (r *SQLBuildRepository) GetByID(ctx context.Context, id int64) (*domain.Build, error) {
	query := `SELECT id, test_suite_id, project_id, build_number, ci_provider, ci_url, created_at, started_at, ended_at, duration, test_case_count FROM builds WHERE id = $1`
	build := &domain.Build{}
	var ciURL sql.NullString
	var startedAt, endedAt sql.NullTime
	var duration sql.NullFloat64
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&build.ID, &build.TestSuiteID, &build.ProjectID, &build.BuildNumber, &build.CIProvider, &ciURL, &build.CreatedAt, &startedAt, &endedAt, &duration, &build.TestCaseCount); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
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
	return build, nil
}

func (r *SQLBuildRepository) GetAllByProjectID(ctx context.Context, projectID int64) ([]*domain.Build, error) {
	query := `SELECT id, test_suite_id, project_id, build_number, ci_provider, ci_url, created_at, started_at, ended_at, duration, test_case_count FROM builds WHERE project_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var builds []*domain.Build
	for rows.Next() {
		build := &domain.Build{}
		var ciURL sql.NullString
		var startedAt, endedAt sql.NullTime
		var duration sql.NullFloat64
		if err := rows.Scan(&build.ID, &build.TestSuiteID, &build.ProjectID, &build.BuildNumber, &build.CIProvider, &ciURL, &build.CreatedAt, &startedAt, &endedAt, &duration, &build.TestCaseCount); err != nil {
			return nil, err
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
		builds = append(builds, build)
	}
	return builds, nil
}

func (r *SQLBuildRepository) GetAllByTestSuiteID(ctx context.Context, suiteID int64) ([]*domain.Build, error) {
	query := `SELECT id, test_suite_id, project_id, build_number, ci_provider, ci_url, created_at, started_at, ended_at, duration, test_case_count FROM builds WHERE test_suite_id = $1 ORDER BY created_at DESC`
	rows, err := r.db.QueryContext(ctx, query, suiteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var builds []*domain.Build
	for rows.Next() {
		build := &domain.Build{}
		var ciURL sql.NullString
		var startedAt, endedAt sql.NullTime
		var duration sql.NullFloat64
		if err := rows.Scan(&build.ID, &build.TestSuiteID, &build.ProjectID, &build.BuildNumber, &build.CIProvider, &ciURL, &build.CreatedAt, &startedAt, &endedAt, &duration, &build.TestCaseCount); err != nil {
			return nil, err
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
		builds = append(builds, build)
	}
	return builds, nil
}

func (r *SQLBuildRepository) Create(ctx context.Context, build *domain.Build) error {
	query := `INSERT INTO builds (test_suite_id, project_id, build_number, ci_provider, ci_url, created_at, started_at, ended_at, duration, test_case_count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`
	var ciURL interface{} = nil
	if build.CIURL != nil {
		ciURL = *build.CIURL
	}
	var startedAt interface{} = nil
	if build.StartedAt != nil {
		startedAt = *build.StartedAt
	}
	var endedAt interface{} = nil
	if build.EndedAt != nil {
		endedAt = *build.EndedAt
	}
	var duration interface{} = nil
	if build.Duration != nil {
		duration = *build.Duration
	}
	return r.db.QueryRowContext(ctx, query, build.TestSuiteID, build.ProjectID, build.BuildNumber, build.CIProvider, ciURL, build.CreatedAt, startedAt, endedAt, duration, build.TestCaseCount).Scan(&build.ID)
}

func (r *SQLBuildRepository) Update(ctx context.Context, id int64, build *domain.Build) (*domain.Build, error) {
	query := `UPDATE builds SET test_suite_id = $1, project_id = $2, build_number = $3, ci_provider = $4, ci_url = $5, started_at = $6, ended_at = $7, duration = $8, test_case_count = $9 WHERE id = $10 RETURNING id, test_suite_id, project_id, build_number, ci_provider, ci_url, created_at, started_at, ended_at, duration, test_case_count`
	var ciURL interface{} = nil
	if build.CIURL != nil {
		ciURL = *build.CIURL
	}
	var startedAt interface{} = nil
	if build.StartedAt != nil {
		startedAt = *build.StartedAt
	}
	var endedAt interface{} = nil
	if build.EndedAt != nil {
		endedAt = *build.EndedAt
	}
	var duration interface{} = nil
	if build.Duration != nil {
		duration = *build.Duration
	}
	updatedBuild := &domain.Build{}
	var ciURLNull sql.NullString
	var startedAtNull, endedAtNull sql.NullTime
	var durationNull sql.NullFloat64
	if err := r.db.QueryRowContext(ctx, query, build.TestSuiteID, build.ProjectID, build.BuildNumber, build.CIProvider, ciURL, startedAt, endedAt, duration, build.TestCaseCount, id).Scan(&updatedBuild.ID, &updatedBuild.TestSuiteID, &updatedBuild.ProjectID, &updatedBuild.BuildNumber, &updatedBuild.CIProvider, &ciURLNull, &updatedBuild.CreatedAt, &startedAtNull, &endedAtNull, &durationNull, &updatedBuild.TestCaseCount); err != nil {
		return nil, err
	}
	if ciURLNull.Valid {
		updatedBuild.CIURL = &ciURLNull.String
	}
	if startedAtNull.Valid {
		updatedBuild.StartedAt = &startedAtNull.Time
	}
	if endedAtNull.Valid {
		updatedBuild.EndedAt = &endedAtNull.Time
	}
	if durationNull.Valid {
		updatedBuild.Duration = &durationNull.Float64
	}
	return updatedBuild, nil
}

func (r *SQLBuildRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM builds WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
