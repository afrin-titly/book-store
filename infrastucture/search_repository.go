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
	baseQuery := "SELECT title, author, genre, price FROM books WHERE 1=1"
	var args []interface{}
	if criteria.Title != "" {
		baseQuery += " AND title LIKE ?"
		args = append(args, "%"+criteria.Title+"%")
	}
	if criteria.Author != "" {
		baseQuery += " AND author LIKE ?"
		args = append(args, "%"+criteria.Author+"%")
	}
	if criteria.Genre != "" {
		baseQuery += " AND genre LIKE ?"
		args = append(args, "%"+criteria.Genre+"%")
	}
	if criteria.Price != "" {
		baseQuery += " AND price LIKE ?"
		args = append(args, "%"+criteria.Price+"%")
	}
	if criteria.SortBy != "" {
		baseQuery += " ORDER BY ?"
		args = append(args, criteria.SortBy)
	}
	if criteria.Order != "" {
		baseQuery += " ?"
		args = append(args, criteria.Order)
	}
	rows, err := query.DB.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
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
