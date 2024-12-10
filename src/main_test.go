package main

import (
	"encoding/json"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	type testCase struct {
		name          string
		route         string
		expecteStatus int
		inputBody     string
		expectedBody  map[string]string
		foundRoute    bool
	}

	tests := []testCase{
		{
			name:          "Create new book",
			route:         "/books/new",
			expecteStatus: 200,
			inputBody:     `{"message": "Hello Gophers!"}`,
			expectedBody:  map[string]string{"message": "Hello Gophers!"},
			foundRoute:    true,
		},
		{
			name:          "Unknown route",
			route:         "/unknown",
			expecteStatus: 404,
			foundRoute:    false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := routes()
			req := httptest.NewRequest("POST", tc.route, strings.NewReader(tc.inputBody))
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if recorder.Code != tc.expecteStatus {
				t.Errorf("expected status %d, got %d", tc.expecteStatus, recorder.Code)
			}

			if tc.foundRoute {
				var body map[string]string
				err := json.Unmarshal(recorder.Body.Bytes(), &body)
				if err != nil {
					t.Fatalf("failed to decode response body: %v", err)
				}

				if !reflect.DeepEqual(body, tc.expectedBody) {
					t.Errorf("expected body %v, got %v", tc.expectedBody, body)
				}
			}
		})
	}
}
