package application_test

import (
	"book-apis/application"
	"book-apis/domain"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockRepository struct {
	Books []domain.Book
	Err   error
}

func (m *MockRepository) GetAll() ([]domain.Book, error) {
	return m.Books, m.Err
}

func TestBookService_GetAll(t *testing.T) {
	type testCase struct {
		name      string
		expected  []domain.Book
		mockSetup *MockRepository
	}
	tests := []testCase{
		{
			name: "Successfull Retrival",
			expected: []domain.Book{
				{Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10},
				{Title: "Test Title 2", Author: "Test Author 2", Genre: "Adventure", Price: "150", Stock: 20},
			},
			mockSetup: &MockRepository{
				Books: []domain.Book{
					{Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10},
					{Title: "Test Title 2", Author: "Test Author 2", Genre: "Adventure", Price: "150", Stock: 20},
				},
				Err: nil,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			service := application.NewBookService(tc.mockSetup)
			result, err := service.GetAll()

			assert.NoError(t, err)
			assert.Equal(t, len(tc.expected), len(result))
		})
	}
}
