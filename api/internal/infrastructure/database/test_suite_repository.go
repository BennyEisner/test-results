package database

import (
	"context"
	"database/sql"

	"github.com/BennyEisner/test-results/internal/domain"
)

type SQLTestSuiteRepository struct {
	db *sql.DB
}

func NewSQLTestSuiteRepository(db *sql.DB) domain.TestSuiteRepository {
	return &SQLTestSuiteRepository{db: db}
}

func (r *SQLTestSuiteRepository) GetByID(ctx context.Context, id int64) (*domain.TestSuite, error) {
	query := `SELECT id, project_id, name, parent_id, time FROM test_suites WHERE id = $1`
	ts := &domain.TestSuite{}
	var parentID sql.NullInt64
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&ts.ID, &ts.ProjectID, &ts.Name, &parentID, &ts.Time); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if parentID.Valid {
		ts.ParentID = &parentID.Int64
	}
	return ts, nil
}

func (r *SQLTestSuiteRepository) GetAllByProjectID(ctx context.Context, projectID int64) ([]*domain.TestSuite, error) {
	query := `SELECT id, project_id, name, parent_id, time FROM test_suites WHERE project_id = $1`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var suites []*domain.TestSuite
	for rows.Next() {
		ts := &domain.TestSuite{}
		var parentID sql.NullInt64
		if err := rows.Scan(&ts.ID, &ts.ProjectID, &ts.Name, &parentID, &ts.Time); err != nil {
			return nil, err
		}
		if parentID.Valid {
			ts.ParentID = &parentID.Int64
		}
		suites = append(suites, ts)
	}
	return suites, nil
}

func (r *SQLTestSuiteRepository) GetByName(ctx context.Context, projectID int64, name string) (*domain.TestSuite, error) {
	query := `SELECT id, project_id, name, parent_id, time FROM test_suites WHERE project_id = $1 AND name = $2`
	ts := &domain.TestSuite{}
	var parentID sql.NullInt64
	if err := r.db.QueryRowContext(ctx, query, projectID, name).Scan(&ts.ID, &ts.ProjectID, &ts.Name, &parentID, &ts.Time); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if parentID.Valid {
		ts.ParentID = &parentID.Int64
	}
	return ts, nil
}

func (r *SQLTestSuiteRepository) Create(ctx context.Context, suite *domain.TestSuite) error {
	query := `INSERT INTO test_suites (project_id, name, parent_id, time) VALUES ($1, $2, $3, $4) RETURNING id`
	var parentID interface{} = nil
	if suite.ParentID != nil {
		parentID = *suite.ParentID
	}
	return r.db.QueryRowContext(ctx, query, suite.ProjectID, suite.Name, parentID, suite.Time).Scan(&suite.ID)
}

func (r *SQLTestSuiteRepository) Update(ctx context.Context, id int64, name string) (*domain.TestSuite, error) {
	query := `UPDATE test_suites SET name = $1 WHERE id = $2 RETURNING id, project_id, name, parent_id, time`
	ts := &domain.TestSuite{}
	var parentID sql.NullInt64
	if err := r.db.QueryRowContext(ctx, query, name, id).Scan(&ts.ID, &ts.ProjectID, &ts.Name, &parentID, &ts.Time); err != nil {
		return nil, err
	}
	if parentID.Valid {
		ts.ParentID = &parentID.Int64
	}
	return ts, nil
}

func (r *SQLTestSuiteRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM test_suites WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
