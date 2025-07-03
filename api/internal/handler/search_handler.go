package handler

import (
	"encoding/json"
	"net/http"

	"github.com/BennyEisner/test-results/internal/models"
	"github.com/BennyEisner/test-results/internal/service"
)

type SearchHandler struct {
	service *service.SearchService
}

func NewSearchHandler(service *service.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

func (h *SearchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	results, err := h.service.Search(query)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode([]models.SearchResult{}); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	if results == nil {
		results = []models.SearchResult{}
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
