package application_test

import (
	"book-apis/application"
	"book-apis/domain"
	"book-apis/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchService_ExecuteSearch(t *testing.T) {
	repo := new(mocks.MockSearchRepository)
	service := application.NewSearchService(repo)
	type testCase struct {
		name      string
		input     domain.SearchCriteria
		expected  []domain.Book
		mockSetup func()
	}
	tests := []testCase{
		{
			name: "Search by title only - found",
			input: domain.SearchCriteria{
				Title: "Harry Potter",
			},
			expected: []domain.Book{
				{Title: "Harry Potter 1", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
				{Title: "Harry Potter 2", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
			},
			mockSetup: func() {
				repo.On("ExecuteSearch", domain.SearchCriteria{Title: "Harry Potter"}).Return([]domain.Book{
					{Title: "Harry Potter 1", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
					{Title: "Harry Potter 2", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
				}, nil)
			},
		},
		{
			name: "Search by title only - not found",
			input: domain.SearchCriteria{
				Title: "Lord of the ring",
			},
			expected: nil,
			mockSetup: func() {
				repo.On("ExecuteSearch", domain.SearchCriteria{Title: "Lord of the ring"}).Return([]domain.Book(nil), nil)
			},
		},
		{
			name: "Search by author only - found",
			input: domain.SearchCriteria{
				Author: "J.K Rowling",
			},
			expected: []domain.Book{
				{Title: "Harry Potter 1", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
				{Title: "Harry Potter 2", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
			},
			mockSetup: func() {
				repo.On("ExecuteSearch", domain.SearchCriteria{Author: "J.K Rowling"}).Return([]domain.Book{
					{Title: "Harry Potter 1", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
					{Title: "Harry Potter 2", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
				}, nil)
			},
		},
		{
			name: "Search by title and author - found",
			input: domain.SearchCriteria{
				Title:  "Harry Potter",
				Author: "J.K Rowling",
			},
			expected: []domain.Book{
				{Title: "Harry Potter 1", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
				{Title: "Harry Potter 2", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
			},
			mockSetup: func() {
				repo.On("ExecuteSearch", domain.SearchCriteria{Title: "Harry Potter", Author: "J.K Rowling"}).Return([]domain.Book{
					{Title: "Harry Potter 1", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
					{Title: "Harry Potter 2", Author: "J.K Rowling", Genre: "Fantasy", Price: "100"},
				}, nil)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, err := service.ExecuteSearch(tc.input)
			if err != nil {
				assert.Error(t, err)
			} else {
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
