package main

import (
	"book-store/src/handlers"
	"net/http"

	"github.com/gorilla/mux"
)

func routes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/books/new", handlers.CreateBookHandler).Methods("POST")
	return r
}

func main() {
	r := routes()
	http.ListenAndServe(":8080", r)
}
