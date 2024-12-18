package interfaces_test

import (
	"book-apis/application"
	"book-apis/domain"
	"book-apis/interfaces"
	"book-apis/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestGetAllBooks(t *testing.T) {
	repo := &mocks.MockBookRepository{}
	service := application.NewBookService(repo)
	h := interfaces.NewBookHandler(service)

	type testCase struct {
		name       string
		expected   []domain.Book
		mockSetup  func()
		statusCode int
	}

	tests := []testCase{
		{
			name: "Successful response",
			expected: []domain.Book{
				{ID: 1, Title: "Book 1", Author: "Author"},
				{ID: 2, Title: "Book 2", Author: "Author"},
			},
			mockSetup: func() {
				repo.On("GetAll").Return([]domain.Book{
					{ID: 1, Title: "Book 1", Author: "Author"},
					{ID: 2, Title: "Book 2", Author: "Author"},
				}, nil).Once()
			},
			statusCode: http.StatusOK,
		},
		{
			name:     "UnSuccessful response",
			expected: nil,
			mockSetup: func() {
				repo.On("GetAll").Return([]domain.Book(nil), errors.New("Some error message")).Once()
			},
			statusCode: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			req, err := http.NewRequest("GET", "/books", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			response := httptest.NewRecorder()
			handler := http.HandlerFunc(h.GetAllBookHandler)
			handler.ServeHTTP(response, req)

			if response.Code != tc.statusCode {
				t.Errorf("Expected status code %d, but got %d", tc.statusCode, response.Code)
			}
			var books []domain.Book
			json.NewDecoder(response.Body).Decode(&books)

			if !reflect.DeepEqual(books, tc.expected) {
				t.Errorf("Handler returned unexpected body:\nGot:  %+v\nWant: %+v", books, tc.expected)
			}
		})
	}

}
