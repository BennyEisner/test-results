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
		SELECT 'project' as type, id, name, '/projects/' || id as url FROM projects WHERE name ILIKE $1
		UNION ALL
		SELECT 'test_suite' as type, id, name, '/suites/' || id as url FROM test_suites WHERE name ILIKE $1
		UNION ALL
		SELECT 'build' as type, id, build_number as name, '/builds/' || id as url FROM builds WHERE build_number ILIKE $1
		UNION ALL
		SELECT 'test_case' as type, id, name, '/cases/' || id as url FROM test_cases WHERE name ILIKE $1
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
