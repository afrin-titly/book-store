package interfaces

import (
	"book-apis/application"
	"book-apis/domain"
	"net/http"
)

type SearchHandler struct {
	service *application.SearchService
}

func NewSearchHandler(service *application.SearchService) *SearchHandler {
	return &SearchHandler{service: service}
}

// GET /books/search?genre=fantasy&sort_by=rating&order=desc
func (s *SearchHandler) ExecuteSearch(w http.ResponseWriter, r *http.Request) {
	var searchCriteria domain.SearchCriteria
	queryParams := r.URL.Query()
	searchCriteria.Title = queryParams.Get("title")
	searchCriteria.Author = queryParams.Get("author")
	searchCriteria.Genre = queryParams.Get("genre")
	searchCriteria.Price = queryParams.Get("price")
	searchCriteria.SortBy = queryParams.Get("sort_by")
	searchCriteria.Order = queryParams.Get("order")

	books, err := s.service.ExecuteSearch(searchCriteria)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
