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

func TestBookService_CreateBook(t *testing.T) {
	mockRepo := new(mocks.MockBookRepository)
	service := application.NewBookService(mockRepo)
	type testCase struct {
		name      string
		input     *domain.Book
		expected  *domain.Book
		mockSetup func()
	}
	tests := []testCase{
		{
			name: "Successful book create",
			input: &domain.Book{
				Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			expected: &domain.Book{
				Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			mockSetup: func() {
				mockRepo.On("CreateBook", &domain.Book{
					Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}).Return(&domain.Book{
					Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}, nil)
			},
		},
		{
			name:     "Un Successful book create",
			input:    &domain.Book{},
			expected: nil,
			mockSetup: func() {
				mockRepo.On("CreateBook", &domain.Book{}).Return(nil, errors.New("Ohh no error!"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, err := service.CreateBook(tc.input)
			if err != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
		mockRepo.AssertExpectations(t)
	}
}

func TestBookService_UpdateBook(t *testing.T) {
	mock := new(mocks.MockBookRepository)
	service := application.NewBookService(mock)
	type testCase struct {
		name      string
		ID        int
		input     *domain.Book
		expected  *domain.Book
		mockSetup func()
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
				mock.On("UpdateBook", &domain.Book{
					Title: "Updated Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}, 1).Return(&domain.Book{
					Title: "Updated Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}, nil)
			},
		},
		{
			name:     "Unsuccessful book update",
			ID:       1,
			input:    nil,
			expected: nil,
			mockSetup: func() {
				mock.On("UpdateBook", (*domain.Book)(nil), 1).Return(nil, errors.New("Oh no error!!"))
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			result, err := service.UpdateBook(tc.input, tc.ID)
			if err != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestBookService_deleteBook(t *testing.T) {
	mock := new(mocks.MockBookRepository)
	service := application.NewBookService(mock)
	type testCase struct {
		name      string
		ID        int
		expected  any
		mockSetup func()
	}

	tests := []testCase{
		{
			name:     "Successfull deletion",
			ID:       1,
			expected: nil,
			mockSetup: func() {
				mock.On("DeleteBook", 1).Return(nil)
			},
		},
		{
			name:     "Unsuccessfull deletion",
			ID:       2,
			expected: "Oh no error!!",
			mockSetup: func() {
				mock.On("DeleteBook", 2).Return(errors.New("Oh no error!"))
			},
		},
	}
	for _, tc := range tests {
		tc.mockSetup()
		err := service.DeleteBook(tc.ID)
		if err != nil {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, err)
		}
	}
}
