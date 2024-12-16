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
	rows, err := r.DB.Query(`SELECT * FROM books`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []domain.Book
	for rows.Next() {
		book := domain.Book{}
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
