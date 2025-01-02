package infrastucture

import (
	"book-apis/domain"
	"database/sql"
	"errors"
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

func (r *BookRepositoryDB) CreateBook(newBook *domain.Book) (*domain.Book, error) {
	book := &domain.Book{}
	err := r.DB.QueryRow(`INSERT INTO books VALUES(?,?,?,?,?)`, newBook.Title, newBook.Author, newBook.Genre, newBook.Price, newBook.Stock).Scan(&book.Title, &book.Author, &book.Genre, &book.Price, &book.Stock)
	if err != nil {
		return nil, err
	}
	return book, err
}

func (r *BookRepositoryDB) UpdateBook(updateBook *domain.Book, ID int) (*domain.Book, error) {
	book := &domain.Book{}
	err := r.DB.QueryRow(`UPDATE books SET title=?, author=?, genre=?, price=?, stock=? WHERE id=?`, updateBook.Title, updateBook.Author, updateBook.Genre, updateBook.Price, updateBook.Stock, ID).Scan(&book.Title, &book.Author, &book.Genre, &book.Price, &book.Stock)
	if err != nil {
		return nil, err
	}
	return book, nil
}

func (r *BookRepositoryDB) DeleteBook(ID int) error {
	result, err := r.DB.Exec(`DELETE FROM books WHERE id=?`, ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("no rows were deleted")
	}
	return nil
}
