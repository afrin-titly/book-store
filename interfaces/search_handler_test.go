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
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchBook(t *testing.T) {
	repo := new(mocks.MockSearchRepository)
	service := application.NewSearchService(repo)
	h := interfaces.NewSearchHandler(service)

	type testCase struct {
		name       string
		expected   []domain.Book
		mockSetup  func()
		statusCode int
		URL        string
	}

	tests := []testCase{
		{
			name: "Successfully return one result",
			expected: []domain.Book{
				{Title: "Beautiful life", Author: "Test Author 1", Genre: "Horror", Price: "100"},
			},
			mockSetup: func() {
				repo.On("ExecuteSearch", domain.SearchCriteria{
					Title: "life",
				}).Return([]domain.Book{
					{Title: "Beautiful life", Author: "Test Author 1", Genre: "Horror", Price: "100"},
				}, nil)
			},
			statusCode: http.StatusOK,
			URL:        "/search?title=life",
		},
		{
			name: "Successfully return multiple result",
			expected: []domain.Book{
				{Title: "Beautiful life", Author: "Test Author 1", Genre: "Horror", Price: "100"},
				{Title: "Sad life", Author: "Test Author 2", Genre: "Reality", Price: "100"},
			},
			mockSetup: func() {
				repo.On("ExecuteSearch", domain.SearchCriteria{
					Title: "life",
				}).Return([]domain.Book{
					{Title: "Beautiful life", Author: "Test Author 1", Genre: "Horror", Price: "100"},
					{Title: "Sad life", Author: "Test Author 2", Genre: "Reality", Price: "100"},
				}, nil)
			},
			statusCode: http.StatusOK,
			URL:        "/search?title=life",
		},
		{
			name: "Search by multiple params",
			expected: []domain.Book{
				{Title: "Beautiful life", Author: "Test Author 1", Genre: "Horror", Price: "100"},
				{Title: "Sad life", Author: "Test Author 2", Genre: "Reality", Price: "100"},
			},
			mockSetup: func() {
				repo.On("ExecuteSearch", domain.SearchCriteria{
					Title:  "life",
					Author: "test",
				}).Return([]domain.Book{
					{Title: "Beautiful life", Author: "Test Author 1", Genre: "Horror", Price: "100"},
					{Title: "Sad life", Author: "Test Author 2", Genre: "Reality", Price: "100"},
				}, nil)
			},
			statusCode: http.StatusOK,
			URL:        "/search?title=life&author=test",
		},
		{
			name:     "Error happens",
			expected: nil,
			mockSetup: func() {
				repo.On("ExecuteSearch", domain.SearchCriteria{
					Title: "books",
				}).Return([]domain.Book(nil), errors.New("oh no error!"))
			},
			statusCode: http.StatusBadRequest,
			URL:        "/search?title=books",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// repo.On mock setup is not being reset between test cases, leading to unexpected behavior in the mock expectations.
			repo.Mock = mock.Mock{}
			tc.mockSetup()
			req, err := http.NewRequest("GET", tc.URL, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			r := mux.NewRouter()
			r.HandleFunc("/search", h.ExecuteSearchHandler).Methods("GET")
			res := httptest.NewRecorder()
			r.ServeHTTP(res, req)

			if res.Code != tc.statusCode {
				t.Errorf("Expected status code %d, but got %d", tc.statusCode, res.Code)
			}

			if res.Code == http.StatusOK {
				var books []domain.Book
				json.NewDecoder(res.Body).Decode(&books)

				assert.Equal(t, tc.expected, books)
			}
		})
		repo.AssertExpectations(t)
	}
}
