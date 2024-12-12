package handlers

import (
	"book-store/src/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestCreateBookHandler(t *testing.T) {
	type testCase struct {
		name             string
		inputRequest     string
		expectedResponse interface{}
		expectedStatus   int
		succussfulReq    bool
	}
	tests := []testCase{
		{
			name:         "Create Book Successfully",
			inputRequest: `{"title":"New book","author":"test author","genre":"horror","price":100,"stock":10}`,
			expectedResponse: &models.Book{
				Title:  "New book",
				Author: "test author",
				Genre:  "horror",
				Price:  100,
				Stock:  10,
			},
			expectedStatus: http.StatusOK,
			succussfulReq:  true,
		},
		{
			name:             "Create Book Unsuccessfull",
			inputRequest:     `Something dummy`,
			expectedResponse: map[string]string{"error": "Malformed request body"},
			expectedStatus:   http.StatusBadRequest,
			succussfulReq:    false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := mux.NewRouter()
			r.HandleFunc("/books", CreateBookHandler).Methods("POST")
			req := httptest.NewRequest("POST", "/books", strings.NewReader(tc.inputRequest))
			recorder := httptest.NewRecorder()
			r.ServeHTTP(recorder, req)

			if recorder.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, got %d", tc.expectedStatus, recorder.Code)
			}

			if tc.succussfulReq {
				var body *models.Book
				json.NewDecoder(recorder.Body).Decode(&body)
				if !reflect.DeepEqual(body, tc.expectedResponse) {
					t.Errorf("Expected body %v, got %v", tc.expectedResponse, body)
				}
			} else {
				var err map[string]string
				json.NewDecoder(recorder.Body).Decode(&err)
				if !reflect.DeepEqual(err, tc.expectedResponse) {
					t.Errorf("Expected body %v, got %v", tc.expectedResponse, err)
				}
			}

		})
	}
}
