package database

import (
	"context"
	"database/sql"

	"github.com/BennyEisner/test-results/internal/build/domain/models"
	"github.com/BennyEisner/test-results/internal/build/domain/ports"
)

type SQLBuildRepository struct {
	db *sql.DB
}

func NewSQLBuildRepository(db *sql.DB) ports.BuildRepository {
	return &SQLBuildRepository{db: db}
}

func (r *SQLBuildRepository) GetBuilds(ctx context.Context, projectID int64, suiteID *int64) ([]*models.Build, error) {
	query := `
		SELECT b.id, ts.project_id, b.test_suite_id, b.build_number, b.duration, b.created_at
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE ts.project_id = $1
	`
	args := []interface{}{projectID}

	if suiteID != nil {
		query += " AND b.test_suite_id = $2"
		args = append(args, *suiteID)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var builds []*models.Build
	for rows.Next() {
		var build models.Build
		var sqlSuiteID sql.NullInt64
		if err := rows.Scan(&build.ID, &build.ProjectID, &sqlSuiteID, &build.BuildNumber, &build.Duration, &build.Timestamp); err != nil {
			return nil, err
		}
		if sqlSuiteID.Valid {
			build.SuiteID = sqlSuiteID.Int64
		}
		builds = append(builds, &build)
	}

	return builds, nil
}

func (r *SQLBuildRepository) GetBuildByID(ctx context.Context, id int64) (*models.Build, error) {
	query := `
		SELECT b.id, ts.project_id, b.test_suite_id, b.build_number, b.duration, b.created_at
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE b.id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var build models.Build
	var sqlSuiteID sql.NullInt64
	if err := row.Scan(&build.ID, &build.ProjectID, &sqlSuiteID, &build.BuildNumber, &build.Duration, &build.Timestamp); err != nil {
		return nil, err
	}
	if sqlSuiteID.Valid {
		build.SuiteID = sqlSuiteID.Int64
	}

	return &build, nil
}

func (r *SQLBuildRepository) CreateBuild(ctx context.Context, build *models.Build) (int64, error) {
	query := "INSERT INTO builds (test_suite_id, build_number, duration, created_at) VALUES ($1, $2, $3, $4) RETURNING id"
	var id int64
	err := r.db.QueryRowContext(ctx, query, build.SuiteID, build.BuildNumber, build.Duration, build.Timestamp).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *SQLBuildRepository) UpdateBuild(ctx context.Context, build *models.Build) error {
	query := "UPDATE builds SET test_suite_id = $1, build_number = $2, duration = $3, created_at = $4 WHERE id = $5"
	_, err := r.db.ExecContext(ctx, query, build.SuiteID, build.BuildNumber, build.Duration, build.Timestamp, build.ID)
	return err
}

func (r *SQLBuildRepository) DeleteBuild(ctx context.Context, id int64) error {
	query := "DELETE FROM builds WHERE id = $1"
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *SQLBuildRepository) GetBuildDurationTrends(ctx context.Context, projectID int64, suiteID int64) ([]*models.BuildDurationTrend, error) {
	query := `
		SELECT b.build_number, b.duration, b.created_at
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE ts.project_id = $1 AND b.test_suite_id = $2
		ORDER BY b.created_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, projectID, suiteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trends []*models.BuildDurationTrend
	for rows.Next() {
		var trend models.BuildDurationTrend
		if err := rows.Scan(&trend.BuildNumber, &trend.Duration, &trend.CreatedAt); err != nil {
			return nil, err
		}
		trends = append(trends, &trend)
	}

	return trends, nil
}

func (r *SQLBuildRepository) GetLatestBuildStatus(ctx context.Context, projectID int64) (string, error) {
	query := `
		SELECT status
		FROM build_test_case_executions
		WHERE build_id IN (
			SELECT id
			FROM builds
			WHERE test_suite_id IN (
				SELECT id
				FROM test_suites
				WHERE project_id = $1
			)
			ORDER BY created_at DESC
			LIMIT 1
		)
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, projectID)

	var status string
	if err := row.Scan(&status); err != nil {
		if err == sql.ErrNoRows {
			return "No builds found", nil
		}
		return "", err
	}

	return status, nil
}

func (r *SQLBuildRepository) GetLatestBuilds(ctx context.Context, projectID int64, limit int) ([]*models.Build, error) {
	query := `
		SELECT b.id, ts.project_id, b.test_suite_id, b.build_number, b.duration, b.created_at
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE ts.project_id = $1
		ORDER BY b.created_at DESC
		LIMIT $2
	`
	rows, err := r.db.QueryContext(ctx, query, projectID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var builds []*models.Build
	for rows.Next() {
		var build models.Build
		var sqlSuiteID sql.NullInt64
		if err := rows.Scan(&build.ID, &build.ProjectID, &sqlSuiteID, &build.BuildNumber, &build.Duration, &build.Timestamp); err != nil {
			return nil, err
		}
		if sqlSuiteID.Valid {
			build.SuiteID = sqlSuiteID.Int64
		}
		builds = append(builds, &build)
	}

	return builds, nil
}
