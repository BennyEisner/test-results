package http

import (
	"encoding/json"
	"net/http"

	"github.com/BennyEisner/test-results/internal/domain/ports"
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
