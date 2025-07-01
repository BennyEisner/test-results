package db

import (
	"database/sql"
	"fmt"

	"github.com/BennyEisner/test-results/internal/models"
)

type SearchRepository struct {
	DB *sql.DB
}

func NewSearchRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{DB: db}
}

func (r *SearchRepository) Search(query string) ([]models.SearchResult, error) {
	rows, err := r.DB.Query(`
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
    `, fmt.Sprintf("%%%s%%", query))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.SearchResult
	for rows.Next() {
		var result models.SearchResult
		if err := rows.Scan(&result.Type, &result.ID, &result.Name, &result.URL); err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
