package infrastucture_test

import (
	"book-apis/domain"
	"book-apis/infrastucture"
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
	}
}
