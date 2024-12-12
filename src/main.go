package main

import (
	"book-store/src/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/books", handlers.CreateBookHandler).Methods("POST")
	r.HandleFunc("/books", handlers.BooksHandler).Methods("GET")
	return r
}

func main() {
	r := routes()
	http.ListenAndServe(":8080", r)
}
