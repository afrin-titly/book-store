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

	"github.com/gorilla/mux"
)

func TestGetAllBooks(t *testing.T) {
	repo := new(mocks.MockBookRepository)
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
				{Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10},
				{Title: "Test Title 2", Author: "Test Author 2", Genre: "Adventure", Price: "150", Stock: 20},
			},
			mockSetup: func() {
				repo.On("GetAll").Return([]domain.Book{
					{Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10},
					{Title: "Test Title 2", Author: "Test Author 2", Genre: "Adventure", Price: "150", Stock: 20},
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

func TestGetOneBook(t *testing.T) {
	type testCase struct {
		name       string
		ID         string
		expected   domain.Book
		mockSetup  func()
		statusCode int
	}

	mockRepo := new(mocks.MockBookRepository)
	service := application.NewBookService(mockRepo)
	h := interfaces.NewBookHandler(service)

	tests := []testCase{
		{
			name: "Successfull - get one Book",
			ID:   "1",
			expected: domain.Book{
				Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			mockSetup: func() {
				mockRepo.On("GetBook", 1).Return(domain.Book{
					Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}, nil).Once()
			},
			statusCode: http.StatusOK,
		},
		{
			name: "Error while converting ID",
			ID:   "abc",
			mockSetup: func() {
				mockRepo.On("GetBook", "abc").Return(nil, errors.New("Conversion Error")).Once()
			},
			statusCode: http.StatusBadRequest,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			req, err := http.NewRequest("GET", "/books/"+tc.ID, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			r := mux.NewRouter()
			r.HandleFunc("/books/{id}", h.GetBookHandler).Methods("GET")
			response := httptest.NewRecorder()
			r.ServeHTTP(response, req)

			if response.Code != tc.statusCode {
				t.Errorf("Expected status code %d, but got %d", tc.statusCode, response.Code)
			}
		})
	}
}
