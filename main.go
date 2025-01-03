package main

import (
	"book-apis/application"
	"book-apis/infrastucture"
	"book-apis/interfaces"
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func routes(h *interfaces.BookHandler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/books", h.GetAllBookHandler).Methods("GET")
	r.HandleFunc("/books/{id}", h.GetBookHandler).Methods("GET")
	r.HandleFunc("/books", h.CreateBookHandler).Methods("POST")
	r.HandleFunc("/books/{id}", h.UpdateBookHandler).Methods("PUT")
	return r
}

func main() {
	connStirng := "host=localhost port=3306 user=mysql password=secret dbname=books sslmode=disable"
	db, err := sql.Open("mysql", connStirng)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	repo := infrastucture.NewBookRepositoryDB(db)
	service := application.NewBookService(repo)
	handler := interfaces.NewBookHandler(service)
	r := routes(handler)
	http.ListenAndServe(":8080", r)
}
