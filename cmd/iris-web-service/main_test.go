package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Elenetta17/iris-web-service/internal/httpapi"
)

// TestServerRoutes tests that routes are properly registered
func TestServerRoutes(t *testing.T) {
	// Create a new ServeMux to simulate main's setup
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", httpapi.FormPage)
	mux.HandleFunc("POST /hello", httpapi.HelloHandler)

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		contentType    string
		expectedStatus int
	}{
		{
			name:           "root GET shows form",
			method:         "GET",
			path:           "/",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "POST /hello handles form submission",
			method:         "POST",
			path:           "/hello",
			body:           "name=Test",
			contentType:    "application/x-www-form-urlencoded",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "GET /hello returns 405",
			method:         "GET",
			path:           "/hello",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "POST / returns 405",
			method:         "POST",
			path:           "/",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "unknown path returns 404",
			method:         "GET",
			path:           "/unknown",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			var err error

			if tt.body != "" {
				req, err = http.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
				if err != nil {
					t.Fatalf("failed to create request: %v", err)
				}
				req.Header.Set("Content-Type", tt.contentType)
			} else {
				req, err = http.NewRequest(tt.method, tt.path, nil)
				if err != nil {
					t.Fatalf("failed to create request: %v", err)
				}
			}

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}
		})
	}
}
