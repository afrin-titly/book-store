package main

import (
	"book-apis/container"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func routes(c *container.Container) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/books", c.BookHandler.GetAllBookHandler).Methods("GET")
	r.HandleFunc("/books/{id}", c.BookHandler.GetBookHandler).Methods("GET")
	r.HandleFunc("/books", c.BookHandler.CreateBookHandler).Methods("POST")
	r.HandleFunc("/books/{id}", c.BookHandler.UpdateBookHandler).Methods("PUT")
	r.HandleFunc("/books/search", c.SearchHandler.ExecuteSearchHandler).Methods("GET")
	return r
}

func main() {
	c := container.NewContainer()
	defer c.DB.Close()

	r := routes(c)
	http.ListenAndServe(":8080", r)
}
