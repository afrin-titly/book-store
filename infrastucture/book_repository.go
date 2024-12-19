package infrastucture

import (
	"book-apis/domain"
	"database/sql"
)

type BookRepositoryDB struct {
	DB *sql.DB
}

func NewBookRepositoryDB(db *sql.DB) *BookRepositoryDB {
	return &BookRepositoryDB{DB: db}
}

func (r *BookRepositoryDB) GetAll() ([]domain.Book, error) {
	rows, err := r.DB.Query(`SELECT title, author, genre, price, stock FROM books`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		book := domain.Book{}
		if err := rows.Scan(&book.Title, &book.Author, &book.Genre, &book.Price, &book.Stock); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}

func (r *BookRepositoryDB) GetBook(ID int) (domain.Book, error) {
	row := r.DB.QueryRow(`SELECT title, author, genre, price, stock FROM books WHERE id = ?`, ID)
	var book domain.Book
	if err := row.Scan(&book.Title, &book.Author, &book.Genre, &book.Price, &book.Stock); err != nil {
		return domain.Book{}, err
	}
	return book, nil
}
