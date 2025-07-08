package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/search/domain/models"
	"github.com/BennyEisner/test-results/internal/search/domain/ports"
)

// SearchService implements the SearchService interface
type SearchService struct {
	repo ports.SearchRepository
}

func NewSearchService(repo ports.SearchRepository) ports.SearchService {
	return &SearchService{repo: repo}
}

func (s *SearchService) Search(ctx context.Context, query string) ([]*models.SearchResult, error) {
	if query == "" {
		return []*models.SearchResult{}, nil
	}

	results, err := s.repo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}
	return results, nil
}
