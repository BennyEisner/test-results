package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
)

type SQLSearchRepository struct {
	db *sql.DB
}

func NewSQLSearchRepository(db *sql.DB) domain.SearchRepository {
	return &SQLSearchRepository{db: db}
}

func (r *SQLSearchRepository) Search(ctx context.Context, query string) ([]*domain.SearchResult, error) {
	searchQuery := `
		SELECT 'project' as type, p.id, p.name, '/projects/' || p.id as url
		FROM projects p
		WHERE p.name ILIKE $1
		UNION ALL
		SELECT 'test_suite' as type, ts.id, ts.name, '/projects/' || ts.project_id || '/suites/' || ts.id as url
		FROM test_suites ts
		WHERE ts.name ILIKE $1
		UNION ALL
		SELECT 'build' as type, b.id, b.build_number as name, '/projects/' || ts.project_id || '/suites/' || b.test_suite_id || '/builds/' || b.id as url
		FROM builds b
		JOIN test_suites ts ON b.test_suite_id = ts.id
		WHERE b.build_number ILIKE $1
	`

	rows, err := r.db.QueryContext(ctx, searchQuery, fmt.Sprintf("%%%s%%", query))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []*domain.SearchResult
	for rows.Next() {
		var result domain.SearchResult
		if err := rows.Scan(&result.Type, &result.ID, &result.Name, &result.URL); err != nil {
			return nil, err
		}
		results = append(results, &result)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
