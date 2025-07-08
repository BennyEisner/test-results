package ports

import (
	"context"

	"github.com/BennyEisner/test-results/internal/search/domain/models"
)

// SearchRepository defines the interface for search data access
type SearchRepository interface {
	Search(ctx context.Context, query string) ([]*models.SearchResult, error)
}

// SearchService defines the interface for search business logic
type SearchService interface {
	Search(ctx context.Context, query string) ([]*models.SearchResult, error)
}
