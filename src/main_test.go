package main

import (
	"net/http/httptest"
	"testing"
)

func TestMain(t *testing.T) {
	type testCase struct {
		name          string
		route         string
		method        string
		expecteStatus int
	}

	tests := []testCase{
		{
			name:          "Create new book",
			route:         "/books/new",
			method:        "POST",
			expecteStatus: 200,
		},
		{
			name:          "Unknown route",
			method:        "GET",
			route:         "/unknown",
			expecteStatus: 404,
		},
		{
			name:          "All books",
			method:        "GET",
			route:         "/books",
			expecteStatus: 200,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := routes()
			req := httptest.NewRequest(tc.method, tc.route, nil)
			recorder := httptest.NewRecorder()
			router.ServeHTTP(recorder, req)

			if recorder.Code != tc.expecteStatus {
				t.Errorf("expected status %d, got %d", tc.expecteStatus, recorder.Code)
			}
		})
	}
}
