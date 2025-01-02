package domain

import "time"

type Book struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Genre     string    `json:"genre"`
	Price     string    `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BookRepository interface {
	GetAll() ([]Book, error)
	GetBook(ID int) (Book, error)
	CreateBook(book *Book) (*Book, error)
	UpdateBook(book *Book, ID int) (*Book, error)
	DeleteBook(ID int) error
}
