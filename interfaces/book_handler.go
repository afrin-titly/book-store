package interfaces

import (
	"book-apis/application"
	"book-apis/domain"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

func (s *BookHandler) GetBookHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Can not convert id to int", http.StatusBadRequest)
		return
	}
	book, err := s.service.GetBook(ID)
	if err != nil {
		http.Error(w, "Can not get Book", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (s *BookHandler) CreateBookHandler(w http.ResponseWriter, r *http.Request) {
	var book domain.Book
	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		http.Error(w, "Can not Decode json", http.StatusBadRequest)
		return
	}
	newBook, err := s.service.CreateBook(&book)
	if err != nil {

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(newBook)
}
