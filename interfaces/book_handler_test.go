package interfaces_test

import (
	"book-apis/domain"
	"book-apis/mocks"
	"testing"
)

func TestBookHandler_GetAll(t *testing.T) {
	type testCase struct {
		name      string
		expected  []domain.Book
		mockSetup func(*mocks.MockBookRepository)
	}

	tests := []testCase{
		{
			name: "Successfull request",
			expected: []domain.Book{
				{ID: 1, Title: "Test Title 1", Author: "Test Author 1"},
				{ID: 2, Title: "Test Title 2", Author: "Test Author 2"},
			},
			mockSetup: func(mockRepo *mocks.MockBookRepository) {
				mockRepo.On("GetAll", []domain.Book{
					{ID: 1, Title: "Book 1", Author: "Author 1"},
					{ID: 2, Title: "Book 2", Author: "Author 2"},
				}).Return([]domain.Book{
					{ID: 1, Title: "Book 1", Author: "Author 1"},
					{ID: 2, Title: "Book 2", Author: "Author 2"},
				}, nil)
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()

		})
	}
}
