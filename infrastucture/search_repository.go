package infrastucture

import (
	"book-apis/domain"
	"database/sql"
	"errors"
)

type SearchRepositoryDB struct {
	DB *sql.DB
}

func NewSearchRepositoryDB(db *sql.DB) *SearchRepositoryDB {
	return &SearchRepositoryDB{DB: db}
}

func (query *SearchRepositoryDB) ExecuteSearch(criteria domain.SearchCriteria) ([]domain.Book, error) {
	queryString := "SELECT title, author, genre, price FROM books WHERE title LIKE ? "
	searchTerm := "%" + criteria.Title + "%"
	rows, err := query.DB.Query(queryString, searchTerm)
	if err != nil {
		return nil, err
	}
	defer query.DB.Close()
	var books []domain.Book
	for rows.Next() {
		var book domain.Book
		if err := rows.Scan(&book.Title, &book.Author, &book.Genre, &book.Price); err != nil {
			return nil, err
		}
		books = append(books, book)
	}

	if len(books) == 0 {
		return nil, errors.New("no result found")
	}
	return books, nil
}
