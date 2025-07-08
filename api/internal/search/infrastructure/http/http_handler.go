package http

import (
	"encoding/json"
	"net/http"

	"github.com/BennyEisner/test-results/internal/search/domain/ports"
)

// SearchHandler handles HTTP requests for search operations
type SearchHandler struct {
	Service ports.SearchService
}

// NewSearchHandler creates a new SearchHandler
func NewSearchHandler(service ports.SearchService) *SearchHandler {
	return &SearchHandler{Service: service}
}

// Search handles GET /search
// @Summary Search across all entities
// @Description Search for projects, test suites, builds, and other entities by name
// @Tags search
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {array} models.SearchResult
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /search [get]
func (h *SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "missing query parameter", http.StatusBadRequest)
		return
	}

	results, err := h.Service.Search(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}
}
