package service

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/BennyEisner/test-results/internal/models"
)

// BuildServiceInterface defines the interface for build service operations.
type BuildServiceInterface interface {
	GetAllBuilds() ([]models.Build, error)
	GetBuildByID(id int64) (*models.Build, error)
	GetRecentBuildsByProjectID(projectID int64) ([]models.Build, error)
	GetBuildsByTestSuiteID(testSuiteID int64) ([]models.Build, error)
	CreateBuild(build *models.Build) (*models.Build, error)
	CreateBuildWithTx(tx *sql.Tx, build *models.Build) (*models.Build, error)
	UpdateBuild(id int64, buildNumber, ciProvider, ciURL *string, duration *float64) (*models.Build, error)
	DeleteBuild(id int64) (int64, error)
	CheckTestSuiteExists(testSuiteID int64) (bool, error)
	GetBuildDurationTrends(projectID, suiteID int64) ([]models.BuildDurationTrend, error)
}

// BuildService provides services related to builds.
type BuildService struct {
	DB *sql.DB
}

// NewBuildService creates a new BuildService.
func NewBuildService(db *sql.DB) *BuildService {
	return &BuildService{DB: db}
}

// CheckTestSuiteExists checks if a test suite with the given ID exists.
func (s *BuildService) CheckTestSuiteExists(testSuiteID int64) (bool, error) {
	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM test_suites WHERE id = $1)", testSuiteID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("database error checking test suite: %w", err)
	}
	return exists, nil
}

// GetAllBuilds fetches all builds from the database.
func (s *BuildService) GetAllBuilds() ([]models.Build, error) {
	const query = `
		SELECT b.id, b.test_suite_id, ts.project_id, b.build_number, b.ci_provider, 
		       b.ci_url, b.created_at, b.test_case_count, b.duration
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		ORDER BY b.created_at DESC
	`
	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("database error fetching all builds: %w", err)
	}
	defer rows.Close()

	var builds []models.Build
	for rows.Next() {
		var b models.Build
		var ciURL sql.NullString
		var duration sql.NullFloat64

		if err := rows.Scan(&b.ID, &b.TestSuiteID, &b.ProjectID, &b.BuildNumber, &b.CIProvider,
			&ciURL, &b.CreatedAt, &b.TestCaseCount, &duration); err != nil {
			return nil, fmt.Errorf("error scanning build: %w", err)
		}

		if ciURL.Valid {
			val := ciURL.String
			b.CIURL = &val
		}

		if duration.Valid {
			val := duration.Float64
			b.Duration = &val
		}

		builds = append(builds, b)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating build rows: %w", err)
	}
	return builds, nil
}

// GetRecentBuildsByProjectID fetches all recent builds for a given projectID.
func (s *BuildService) GetRecentBuildsByProjectID(projectID int64) ([]models.Build, error) {
	const query = `
		SELECT b.id, b.test_suite_id, ts.project_id, b.build_number, b.ci_provider, 
		       b.ci_url, b.created_at, b.test_case_count, b.duration
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE ts.project_id = $1
		ORDER BY b.created_at DESC
	`
	rows, err := s.DB.Query(query, projectID)
	if err != nil {
		return nil, fmt.Errorf("database error fetching builds for project ID %d: %w", projectID, err)
	}
	defer rows.Close()

	var builds []models.Build
	for rows.Next() {
		var b models.Build
		var ciURL sql.NullString
		var duration sql.NullFloat64

		if err := rows.Scan(&b.ID, &b.TestSuiteID, &b.ProjectID, &b.BuildNumber, &b.CIProvider,
			&ciURL, &b.CreatedAt, &b.TestCaseCount, &duration); err != nil {
			return nil, fmt.Errorf("error scanning build for project ID %d: %w", projectID, err)
		}

		if ciURL.Valid {
			val := ciURL.String
			b.CIURL = &val
		}

		if duration.Valid {
			val := duration.Float64
			b.Duration = &val
		}

		builds = append(builds, b)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating build rows for project ID %d: %w", projectID, err)
	}
	return builds, nil
}

// GetBuildByID fetches a single build by its ID.
func (s *BuildService) GetBuildByID(id int64) (*models.Build, error) {
	const query = `
		SELECT b.id, b.test_suite_id, ts.project_id, b.build_number, b.ci_provider, 
		       b.ci_url, b.created_at, b.test_case_count, b.duration
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE b.id = $1
	`
	var b models.Build
	var ciURL sql.NullString
	var duration sql.NullFloat64

	err := s.DB.QueryRow(query, id).Scan(
		&b.ID, &b.TestSuiteID, &b.ProjectID, &b.BuildNumber, &b.CIProvider,
		&ciURL, &b.CreatedAt, &b.TestCaseCount, &duration)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // Let handler decide on 404 or other error
		}
		return nil, fmt.Errorf("database error fetching build: %w", err)
	}

	if ciURL.Valid {
		val := ciURL.String
		b.CIURL = &val
	}

	if duration.Valid {
		val := duration.Float64
		b.Duration = &val
	}

	return &b, nil
}

// GetBuildsByTestSuiteID fetches all builds for a given testSuiteID.
func (s *BuildService) GetBuildsByTestSuiteID(testSuiteID int64) ([]models.Build, error) {
	const query = `
		SELECT b.id, b.test_suite_id, ts.project_id, b.build_number, b.ci_provider, 
		       b.ci_url, b.created_at, b.test_case_count, b.duration
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE b.test_suite_id = $1
		ORDER BY b.created_at DESC
	`
	rows, err := s.DB.Query(query, testSuiteID)
	if err != nil {
		return nil, fmt.Errorf("database error fetching builds for test suite ID %d: %w", testSuiteID, err)
	}
	defer rows.Close()

	var builds []models.Build
	for rows.Next() {
		var b models.Build
		var ciURL sql.NullString
		var duration sql.NullFloat64

		if err := rows.Scan(&b.ID, &b.TestSuiteID, &b.ProjectID, &b.BuildNumber, &b.CIProvider,
			&ciURL, &b.CreatedAt, &b.TestCaseCount, &duration); err != nil {
			return nil, fmt.Errorf("error scanning build for test suite ID %d: %w", testSuiteID, err)
		}

		if ciURL.Valid {
			val := ciURL.String
			b.CIURL = &val
		}

		if duration.Valid {
			val := duration.Float64
			b.Duration = &val
		}

		builds = append(builds, b)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating build rows for test suite ID %d: %w", testSuiteID, err)
	}
	return builds, nil
}

// CreateBuild creates a new build in the database.
// The input models.Build should have TestSuiteID, BuildNumber, CIProvider, and optionally CIURL populated.
// CreatedAt will be set by the database (NOW()). ID will be set from RETURNING.
func (s *BuildService) CreateBuild(build *models.Build) (*models.Build, error) {
	var newBuildID int64
	var createdAt time.Time

	var ciURLNullStr sql.NullString
	if build.CIURL != nil && strings.TrimSpace(*build.CIURL) != "" {
		ciURLNullStr = sql.NullString{String: *build.CIURL, Valid: true}
	} else {
		ciURLNullStr = sql.NullString{String: "", Valid: false}
	}

	var durationNull sql.NullFloat64
	if build.Duration != nil {
		durationNull = sql.NullFloat64{Float64: *build.Duration, Valid: true}
	} else {
		durationNull = sql.NullFloat64{Float64: 0, Valid: false}
	}

	err := s.DB.QueryRow(
		"INSERT INTO builds(test_suite_id, build_number, ci_provider, ci_url, created_at, test_case_count, duration) VALUES($1, $2, $3, $4, NOW(), $5, $6) RETURNING id, created_at",
		build.TestSuiteID, build.BuildNumber, build.CIProvider, ciURLNullStr, build.TestCaseCount, durationNull,
	).Scan(&newBuildID, &createdAt)

	if err != nil {
		return nil, fmt.Errorf("database error creating build: %w", err)
	}

	// Fetch the created build with project_id
	return s.GetBuildByID(newBuildID)
}

// CreateBuildWithTx creates a new build in the database within an existing transaction.
// The input models.Build should have TestSuiteID, BuildNumber, CIProvider, and optionally CIURL populated.
// CreatedAt will be set by the database (NOW()). ID will be set from RETURNING.
func (s *BuildService) CreateBuildWithTx(tx *sql.Tx, build *models.Build) (*models.Build, error) {
	var newBuildID int64
	var createdAt time.Time

	var ciURLNullStr sql.NullString
	if build.CIURL != nil && strings.TrimSpace(*build.CIURL) != "" {
		ciURLNullStr = sql.NullString{String: *build.CIURL, Valid: true}
	} else {
		ciURLNullStr = sql.NullString{String: "", Valid: false}
	}

	var durationNull sql.NullFloat64
	if build.Duration != nil {
		durationNull = sql.NullFloat64{Float64: *build.Duration, Valid: true}
	} else {
		durationNull = sql.NullFloat64{Float64: 0, Valid: false}
	}

	// Use tx.QueryRow instead of s.DB.QueryRow
	err := tx.QueryRow(
		"INSERT INTO builds(test_suite_id, build_number, ci_provider, ci_url, created_at, test_case_count, duration) VALUES($1, $2, $3, $4, NOW(), $5, $6) RETURNING id, created_at",
		build.TestSuiteID, build.BuildNumber, build.CIProvider, ciURLNullStr, build.TestCaseCount, durationNull,
	).Scan(&newBuildID, &createdAt)

	if err != nil {
		return nil, fmt.Errorf("database error creating build with tx: %w", err)
	}

	createdBuild := &models.Build{
		ID:            newBuildID,
		TestSuiteID:   build.TestSuiteID,
		BuildNumber:   build.BuildNumber,
		CIProvider:    build.CIProvider,
		CIURL:         build.CIURL,
		CreatedAt:     createdAt,
		TestCaseCount: build.TestCaseCount,
		Duration:      build.Duration,
	}
	return createdBuild, nil
}

// DeleteBuild deletes a build by its ID and returns the number of rows affected.
func (s *BuildService) DeleteBuild(id int64) (int64, error) {
	result, err := s.DB.Exec("DELETE FROM builds WHERE id = $1", id)
	if err != nil {
		return 0, fmt.Errorf("database error deleting build: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("error checking delete result: %w", err)
	}
	return rowsAffected, nil
}

// UpdateBuild updates an existing build.
// Only non-nil fields in the input will be updated.
func (s *BuildService) UpdateBuild(id int64, buildNumber, ciProvider, ciURL *string, duration *float64) (*models.Build, error) {
	// First, check if build exists
	_, err := s.GetBuildByID(id) // Leverage existing method
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err // Propagate ErrNoRows for handler to interpret as 404
		}
		return nil, fmt.Errorf("error checking build existence before update: %w", err)
	}

	updateFields := []string{}
	args := []interface{}{}
	argID := 1

	if buildNumber != nil {
		if strings.TrimSpace(*buildNumber) == "" {
			return nil, fmt.Errorf("build number cannot be empty if provided for update")
		}
		updateFields = append(updateFields, fmt.Sprintf("build_number = $%d", argID))
		args = append(args, *buildNumber)
		argID++
	}

	if ciProvider != nil {
		if strings.TrimSpace(*ciProvider) == "" {
			return nil, fmt.Errorf("ci provider cannot be empty if provided for update")
		}
		updateFields = append(updateFields, fmt.Sprintf("ci_provider = $%d", argID))
		args = append(args, *ciProvider)
		argID++
	}

	if ciURL != nil { // If CIURL key is present
		ciURLToUpdate := sql.NullString{String: *ciURL, Valid: strings.TrimSpace(*ciURL) != ""}
		updateFields = append(updateFields, fmt.Sprintf("ci_url = $%d", argID))
		args = append(args, ciURLToUpdate)
		argID++
	}

	if duration != nil {
		updateFields = append(updateFields, fmt.Sprintf("duration = $%d", argID))
		args = append(args, sql.NullFloat64{Float64: *duration, Valid: true})
		argID++
	}

	if len(updateFields) == 0 {
		// No fields to update, maybe return current build or an error
		// For now, let's return the current build as if no update occurred, or an error
		return nil, fmt.Errorf("no valid fields provided for update")
	}

	args = append(args, id)
	query := fmt.Sprintf("UPDATE builds SET %s WHERE id = $%d",
		strings.Join(updateFields, ", "), argID)

	_, err = s.DB.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("update build failed: %w", err)
	}

	// After successful update, fetch the updated build with project_id
	return s.GetBuildByID(id)
}

// GetBuildDurationTrends fetches build duration trends for a given project ID.
func (s *BuildService) GetBuildDurationTrends(projectID, suiteID int64) ([]models.BuildDurationTrend, error) {
	const query = `
		SELECT b.build_number, b.duration, b.created_at
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE ts.project_id = $1 AND b.test_suite_id = $2 AND b.duration IS NOT NULL
		ORDER BY b.created_at DESC
		LIMIT 50
	`

	rows, err := s.DB.Query(query, projectID, suiteID)
	if err != nil {
		return nil, fmt.Errorf("database error fetching build duration trends for project ID %d: %w", projectID, err)
	}
	defer rows.Close()

	var trends []models.BuildDurationTrend
	for rows.Next() {
		var trend models.BuildDurationTrend
		if err := rows.Scan(&trend.BuildNumber, &trend.Duration, &trend.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning build duration trend for project ID %d: %w", projectID, err)
		}
		trends = append(trends, trend)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating build duration trend rows for project ID %d: %w", projectID, err)
	}

	return trends, nil
}
