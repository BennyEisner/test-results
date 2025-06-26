// Uses the repository to make search operations

package service

import (
	"database/sql"

	"github.com/BennyEisner/test-results/internal/db"
	"github.com/BennyEisner/test-results/internal/models"
)

type SearchService struct {
	repo *db.SearchRepository
}

func NewSearchService(database *sql.DB) *SearchService {
	return &SearchService{
		repo: db.NewSearchRepository(database),
	}
}

func (s *SearchService) Search(query string) ([]models.SearchResult, error) {
	if query == "" {
		return []models.SearchResult{}, nil
	}

	return s.repo.Search(query)
}
