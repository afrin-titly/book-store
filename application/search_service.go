package application

import (
	"book-apis/domain"
)

type SearchService struct {
	search domain.SearchRepository
}

func NewSearchService(search domain.SearchRepository) *SearchService {
	return &SearchService{search: search}
}

func (s *SearchService) ExecuteSearch(parameters domain.SearchCriteria) ([]domain.Book, error) {
	return s.search.ExecuteSearch(parameters)
}
