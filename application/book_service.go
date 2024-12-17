package application

import "book-apis/domain"

type BookService struct {
	service domain.BookRepository
}

func NewBookService(repo domain.BookRepository) *BookService {
	return &BookService{service: repo}
}

func (s *BookService) GetAll() ([]domain.Book, error) {
	return s.service.GetAll()
}
