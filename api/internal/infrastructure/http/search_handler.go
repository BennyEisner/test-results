package http

import (
	"net/http"

	"github.com/BennyEisner/test-results/internal/domain"
)

type SearchHandler struct {
	service domain.SearchService
}

func NewSearchHandler(service domain.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	results, err := h.service.Search(r.Context(), query)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "failed to perform search")
		return
	}

	if results == nil {
		results = []*domain.SearchResult{}
	}

	respondWithJSON(w, http.StatusOK, results)
}
