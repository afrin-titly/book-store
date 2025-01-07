package infrastucture_test

import (
	"book-apis/domain"
	"book-apis/infrastucture"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestSearchRepositoryDB_ExecuteSearch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error initializing sqlmock: %v", err)
	}
	defer db.Close()
	type testCase struct {
		name      string
		input     domain.SearchCriteria
		expected  any
		mockSetup func()
	}
	tests := []testCase{
		{
			name: "Search by title - success",
			input: domain.SearchCriteria{
				Title: "Harry Potter",
			},
			expected: []domain.Book{
				{Title: "Harry Potter 1", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
				{Title: "Harry Potter 2", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"title", "author", "genre", "price"}).AddRow("Harry Potter 1", "J.K Rowling", "Fantasy", "100").AddRow("Harry Potter 2", "J.K Rowling", "Fantasy", "100")
				mock.ExpectQuery(`^SELECT title, author, genre, price FROM books WHERE 1=1 AND title LIKE \?$`).WithArgs("%Harry Potter%").WillReturnRows(rows)
			},
		},
		{
			name: "Search by title - no result",
			input: domain.SearchCriteria{
				Title: "aest blah",
			},
			expected: "no result found",
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"title", "author", "genre", "price"})
				mock.ExpectQuery(`^SELECT title, author, genre, price FROM books WHERE 1=1 AND title LIKE \?$`).WithArgs(`%aest blah%`).WillReturnRows(rows)
			},
		},
		{
			name: "Search by title and author - success",
			input: domain.SearchCriteria{
				Title:  "Harry",
				Author: "Rowling",
			},
			expected: []domain.Book{
				{Title: "Harry Potter 1", Author: "G. Rowling", Genre: "Fantasy", Price: "100"},
				{Title: "Harry James 2", Author: "J.K Rowling", Genre: "Horror", Price: "100"},
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"title", "author", "genre", "price"}).AddRow("Harry Potter 1", "G. Rowling", "Fantasy", "100").AddRow("Harry James 2", "J.K Rowling", "Horror", "100")
				mock.ExpectQuery(`^SELECT title, author, genre, price FROM books WHERE 1=1 AND title LIKE \? AND author LIKE \?$`).
					WithArgs("%Harry%", "%Rowling%").
					WillReturnRows(rows)
			},
		},
		{
			name: "Search by title and Sort by author - success",
			input: domain.SearchCriteria{
				Title:  "Harry",
				SortBy: "author",
			},
			expected: []domain.Book{
				{Title: "Harry Potter 2", Author: "A. Rowling", Genre: "Fantasy", Price: "100"},
				{Title: "Harry James 1", Author: "J.K Rowling", Genre: "Horror", Price: "100"},
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"title", "author", "genre", "price"}).AddRow("Harry Potter 2", "A. Rowling", "Fantasy", "100").AddRow("Harry James 1", "J.K Rowling", "Horror", "100")
				mock.ExpectQuery(`^SELECT title, author, genre, price FROM books WHERE 1=1 AND title LIKE \? ORDER BY \?`).
					WithArgs("%Harry%", "author").
					WillReturnRows(rows)
			},
		},
		{
			name: "Search by title and Sort by author ASC - success",
			input: domain.SearchCriteria{
				Title:  "Harry",
				SortBy: "author",
				Order:  "asc",
			},
			expected: []domain.Book{
				{Title: "Harry Potter 2", Author: "A. Rowling", Genre: "Fantasy", Price: "100"},
				{Title: "Harry James 1", Author: "J.K Rowling", Genre: "Horror", Price: "100"},
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"title", "author", "genre", "price"}).AddRow("Harry Potter 2", "A. Rowling", "Fantasy", "100").AddRow("Harry James 1", "J.K Rowling", "Horror", "100")
				mock.ExpectQuery(`^SELECT title, author, genre, price FROM books WHERE 1=1 AND title LIKE \? ORDER BY \? \?`).
					WithArgs("%Harry%", "author", "asc").
					WillReturnRows(rows)
			},
		},
		{
			name: "Search by title and Sort by price DESC - success",
			input: domain.SearchCriteria{
				Title:  "Harry",
				SortBy: "price",
				Order:  "desc",
			},
			expected: []domain.Book{
				{Title: "Harry James 1", Author: "J.K Rowling", Genre: "Horror", Price: "100"},
				{Title: "Harry Potter 2", Author: "J.K Rowling", Genre: "Fantasy", Price: "50"},
			},
			mockSetup: func() {
				rows := sqlmock.NewRows([]string{"title", "author", "genre", "price"}).AddRow("Harry James 1", "J.K Rowling", "Horror", "100").AddRow("Harry Potter 2", "J.K Rowling", "Fantasy", "50")
				mock.ExpectQuery(`^SELECT title, author, genre, price FROM books WHERE 1=1 AND title LIKE \? ORDER BY \? \?$`).
					WithArgs("%Harry%", "price", "desc").WillReturnRows(rows)
			},
		},
	}
	repo := infrastucture.NewSearchRepositoryDB(db)
	for _, tc := range tests {
		tc.mockSetup()
		result, err := repo.ExecuteSearch(tc.input)
		if err != nil {
			assert.Error(t, err)
			assert.Nil(t, result)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		}
	}
	assert.NoError(t, mock.ExpectationsWereMet())
}
