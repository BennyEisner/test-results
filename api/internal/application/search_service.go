package application

import (
	"context"
	"fmt"

	"github.com/BennyEisner/test-results/internal/domain"
)

type SearchService struct {
	searchRepo domain.SearchRepository
}

func NewSearchService(searchRepo domain.SearchRepository) domain.SearchService {
	return &SearchService{
		searchRepo: searchRepo,
	}
}

func (s *SearchService) Search(ctx context.Context, query string) ([]*domain.SearchResult, error) {
	if query == "" {
		return []*domain.SearchResult{}, nil
	}

	results, err := s.searchRepo.Search(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to search: %w", err)
	}

	return results, nil
}
