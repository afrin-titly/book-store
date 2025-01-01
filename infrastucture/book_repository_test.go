package infrastucture_test

import (
	"book-apis/domain"
	"book-apis/infrastucture"
	"database/sql"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestBookRepositoryDB_GetAll(t *testing.T) {
	type testCase struct {
		name        string
		expected    []domain.Book
		mockSetup   func()
		shouldError bool
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %v", err)
	}

	defer db.Close()

	repo := infrastucture.NewBookRepositoryDB(db)

	tests := []testCase{
		{
			name: "success - fetch all books",
			expected: []domain.Book{
				{Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10},
				{Title: "Test Title 2", Author: "Test Author 2", Genre: "Adventure", Price: "150", Stock: 20},
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"title", "author", "genre", "price", "stock"}).AddRow("Test Title 1", "Test Author 1", "Horror", "100", 10).AddRow("Test Title 2", "Test Author 2", "Adventure", "150", 20)
				mock.ExpectQuery("SELECT title, author, genre, price, stock FROM books").WillReturnRows(rows)
			},
			shouldError: false,
		},
		{
			name:     "failure - query execution fails",
			expected: nil,
			mockSetup: func() {
				mock.ExpectQuery("SELECT title, author, genre, price, stock FROM books").WillReturnError(fmt.Errorf("Some DB error"))
			},
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

			books, err := repo.GetAll()

			if tc.shouldError {
				assert.Error(t, err)
				assert.Nil(t, books)
			} else {
				assert.NoError(t, err)
				assert.Len(t, books, 2)
				assert.Equal(t, tc.expected, books)
			}
		})
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}
}

func TestBookRepositoryDB_GetOneBook(t *testing.T) {
	type testCase struct {
		name        string
		ID          int
		expected    domain.Book
		mockSetup   func()
		shouldError bool
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %v", err)
	}

	defer db.Close()

	repo := infrastucture.NewBookRepositoryDB(db)

	tests := []testCase{
		{
			name: "success - fetch one book",
			ID:   1,
			expected: domain.Book{
				Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			mockSetup: func() {
				row := sqlmock.NewRows([]string{"title", "author", "genre", "price", "stock"}).AddRow("Test Title 1", "Test Author 1", "Horror", "100", 10)
				mock.ExpectQuery("SELECT title, author, genre, price, stock FROM books WHERE id = ?").WithArgs(1).WillReturnRows(row)
			},
			shouldError: false,
		},
		{
			name:     "not success - fetch one book",
			ID:       2,
			expected: domain.Book{},
			mockSetup: func() {
				mock.ExpectQuery("SELECT title, author, genre, price, stock FROM books WHERE id = ?").WithArgs(2).WillReturnError(sql.ErrNoRows)
			},
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			book, err := repo.GetBook(tc.ID)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Equal(t, tc.expected, domain.Book{})
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, book)
			}
		})
		err = mock.ExpectationsWereMet()
		assert.NoError(t, err)
	}
}

func TestBookRepositoryDB_CreateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %v", err)
	}

	defer db.Close()
	type testCase struct {
		name        string
		input       *domain.Book
		expected    *domain.Book
		mockSetup   func()
		shouldError bool
	}

	tests := []testCase{
		{
			name: "success - create book",
			input: &domain.Book{
				Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			expected: &domain.Book{
				Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			mockSetup: func() {
				row := sqlmock.NewRows([]string{"title", "author", "genre", "price", "stock"}).AddRow("Test Title 1", "Test Author 1", "Horror", "100", 10)
				mock.ExpectQuery("INSERT INTO books").WithArgs("Test Title 1", "Test Author 1", "Horror", "100", 10).WillReturnRows(row)
			},
			shouldError: false,
		},
		{
			name: "not success - create book",
			input: &domain.Book{
				Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			mockSetup: func() {
				mock.ExpectQuery("INSERT INTO books").WithArgs().WillReturnError(fmt.Errorf("Ohh no! Error!"))
			},
			shouldError: true,
		},
	}
	repo := infrastucture.NewBookRepositoryDB(db)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, err := repo.CreateBook(tc.input)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}

func TestBookRepositoryDB_UpdateBook(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %v", err)
	}
	defer db.Close()
	type testCase struct {
		name        string
		ID          int
		input       *domain.Book
		expected    *domain.Book
		mockSetup   func()
		shouldError bool
	}
	tests := []testCase{
		{
			name: "Successful book update",
			ID:   1,
			input: &domain.Book{
				Title: "Updated Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			expected: &domain.Book{
				Title: "Updated Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			mockSetup: func() {
				row := sqlmock.NewRows([]string{"title", "author", "genre", "price", "stock"}).AddRow("Updated Test Title 1", "Test Author 1", "Horror", "100", 10)
				mock.ExpectQuery("UPDATE books").WithArgs("Updated Test Title 1", "Test Author 1", "Horror", "100", 10, 1).WillReturnRows(row)
			},
			shouldError: false,
		},
		{
			name:     "Unsuccessful book update",
			ID:       1,
			input:    &domain.Book{},
			expected: nil,
			mockSetup: func() {
				mock.ExpectQuery("UPDATE books").WithArgs().WillReturnError(fmt.Errorf("Oh no error!!"))
			},
			shouldError: true,
		},
	}
	repo := infrastucture.NewBookRepositoryDB(db)
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, err := repo.UpdateBook(tc.input, tc.ID)
			if tc.shouldError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
