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
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/mock"
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

func TestCreateBook(t *testing.T) {
	mock := new(mocks.MockBookRepository)
	service := application.NewBookService(mock)
	h := interfaces.NewBookHandler(service)
	type testCase struct {
		name        string
		input       string
		expected    domain.Book
		mockSetup   func()
		statusCode  int
		shouldError bool
	}

	tests := []testCase{
		{
			name:  "Succussful book create",
			input: `{"title": "Test Title 1", "author": "Test Author 1", "genre": "Horror", "price": "100", "stock": 10}`,
			expected: domain.Book{
				Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			mockSetup: func() {
				mock.On("CreateBook", &domain.Book{
					Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}).Return(
					&domain.Book{
						Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
					}, nil)
			},
			statusCode:  http.StatusOK,
			shouldError: false,
		},
		{
			name:     "json decode fail",
			input:    `{title: "Test Title 1", "author": "Test Author 1", "genre": "Horror", "price": "100", "stock": 10}`,
			expected: domain.Book{},
			mockSetup: func() {
				mock.On("CreateBook", &domain.Book{
					Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}).Return(
					&domain.Book{
						Title: "Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
					}, nil)
			},
			statusCode:  http.StatusBadRequest,
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			req, err := http.NewRequest("POST", "/books", strings.NewReader(tc.input))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			r := mux.NewRouter()
			r.HandleFunc("/books", h.CreateBookHandler).Methods("POST")
			response := httptest.NewRecorder()
			r.ServeHTTP(response, req)

			if tc.statusCode != response.Code {
				t.Errorf("Expected status code %d, but got %d", tc.statusCode, response.Code)
			}

			if !tc.shouldError {
				var newBook domain.Book
				json.NewDecoder(response.Body).Decode(&newBook)
				if newBook != tc.expected {
					t.Errorf("Expected body %v, but got %v", tc.expected, newBook)
				}
			}
		})
	}
}

func TestUpdateBook(t *testing.T) {
	repo := new(mocks.MockBookRepository)
	service := application.NewBookService(repo)
	h := interfaces.NewBookHandler(service)
	type testCase struct {
		name        string
		ID          string
		input       string
		expected    domain.Book
		mockSetup   func()
		statusCode  int
		shouldError bool
	}
	tests := []testCase{
		{
			name:  "Successfully update book",
			ID:    "1",
			input: `{"title": "Updated Test Title 1", "author": "Test Author 1", "genre": "Horror", "price": "100", "stock": 10}`,
			expected: domain.Book{
				Title: "Updated Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
			},
			mockSetup: func() {
				repo.On("UpdateBook", &domain.Book{
					Title: "Updated Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}, 1).Return(&domain.Book{
					Title: "Updated Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}, nil)
			},
			statusCode:  http.StatusOK,
			shouldError: false,
		},
		{
			name:     "json decode fail",
			ID:       "1",
			input:    `{title: "Updated Test Title 1", "author": "Test Author 1", "genre": "Horror", "price": "100", "stock": 10}`,
			expected: domain.Book{},
			mockSetup: func() {
				repo.On("UpdateBook", &domain.Book{
					Title: "Updated Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}, 1).Return(&domain.Book{
					Title: "Updated Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}, nil)
			},
			statusCode:  http.StatusBadRequest,
			shouldError: true,
		},
		{
			name:     "ID retirval failed",
			ID:       "a",
			input:    `{"title": "Updated Test Title 1", "author": "Test Author 1", "genre": "Horror", "price": "100", "stock": 10}`,
			expected: domain.Book{},
			mockSetup: func() {
				repo.On("UpdateBook", &domain.Book{
					Title: "Updated Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}, 1).Return(&domain.Book{
					Title: "Updated Test Title 1", Author: "Test Author 1", Genre: "Horror", Price: "100", Stock: 10,
				}, nil)
			},
			statusCode:  http.StatusBadRequest,
			shouldError: true,
		},
		{
			name:     "update failed",
			ID:       "10",
			input:    `{"title": "Updated Test Title 1", "author": "Test Author 1", "genre": "Horror", "price": "100", "stock": 10}`,
			expected: domain.Book{},
			mockSetup: func() {
				repo.On("UpdateBook", mock.AnythingOfType("*domain.Book"), 10).Return(nil, errors.New("Oh no error!"))
			},
			statusCode:  http.StatusBadRequest,
			shouldError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			req, err := http.NewRequest("PUT", "/books/"+tc.ID, strings.NewReader(tc.input))
			if err != nil {
				t.Errorf("Failed to create request %v", err)
			}
			r := mux.NewRouter()
			r.HandleFunc("/books/{id}", h.UpdateBookHandler).Methods("PUT")
			response := httptest.NewRecorder()
			r.ServeHTTP(response, req)

			if tc.statusCode != response.Code {
				t.Errorf("Expected status code %d, but got %d", tc.statusCode, response.Code)
			}
			if !tc.shouldError {
				var updatedBook domain.Book
				json.NewDecoder(response.Body).Decode(&updatedBook)
				if updatedBook != tc.expected {
					t.Errorf("Expected body %v, but got %v", tc.expected, updatedBook)
				}
			}
		})
	}
}
