package interfaces

import (
	"book-apis/application"
	"encoding/json"
	"net/http"
)

type BookHandler struct {
	service *application.BookService
}

func NewBookHandler(service *application.BookService) *BookHandler {
	return &BookHandler{service: service}
}

func (s *BookHandler) GetAllBookHandler(w http.ResponseWriter, r *http.Request) {
	books, err := s.service.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(books)
}
