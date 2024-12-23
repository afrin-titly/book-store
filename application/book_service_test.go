package application_test

import (
	"book-apis/application"
	"book-apis/domain"
	"book-apis/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookService_GetAll(t *testing.T) {
	type testCase struct {
		name      string
		expected  []domain.Book
		mockSetup func(mockRepo *mocks.MockBookRepository)
	}
	tests := []testCase{
		{
			name: "Successful Retrieval",
			expected: []domain.Book{
				{Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10},
				{Title: "Test Title 2", Author: "Test Author 2", Genre: "Adventure", Price: "150", Stock: 20},
			},
			mockSetup: func(mockRepo *mocks.MockBookRepository) {
				mockRepo.On("GetAll").Return([]domain.Book{
					{Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10},
					{Title: "Test Title 2", Author: "Test Author 2", Genre: "Adventure", Price: "150", Stock: 20},
				}, nil)
			},
		},
		{
			name:     "Empty List",
			expected: []domain.Book{},
			mockSetup: func(mockRepo *mocks.MockBookRepository) {
				mockRepo.On("GetAll").Return([]domain.Book{}, nil)
			},
		},
		{
			name:     "Error Retrieval",
			expected: nil,
			mockSetup: func(mockRepo *mocks.MockBookRepository) {
				mockRepo.On("GetAll").Return([]domain.Book(nil), assert.AnError)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockBookRepository)

			tc.mockSetup(mockRepo)

			service := application.NewBookService(mockRepo)

			result, err := service.GetAll()

			if tc.expected != nil {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			} else {
				assert.Error(t, err)
				assert.Nil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestBookService_GetOneBook(t *testing.T) {
	type testCase struct {
		name      string
		ID        int
		expected  domain.Book
		mockSetup func(mockRepo *mocks.MockBookRepository)
	}
	tests := []testCase{
		{
			name: "Successful retrival one book",
			ID:   1,
			expected: domain.Book{
				Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			mockSetup: func(mockRepo *mocks.MockBookRepository) {
				mockRepo.On("GetBook", 1).Return(domain.Book{
					Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}, nil)
			},
		},
		{
			name:     "Not Successful retrival one book",
			ID:       2,
			expected: domain.Book{},
			mockSetup: func(mockRepo *mocks.MockBookRepository) {
				mockRepo.On("GetBook", 2).Return(domain.Book{}, errors.New("Some error happened"))
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockBookRepository)
			tc.mockSetup(mockRepo)
			service := application.NewBookService(mockRepo)
			result, err := service.GetBook(tc.ID)
			if tc.expected != (domain.Book{}) {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tc.expected, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
