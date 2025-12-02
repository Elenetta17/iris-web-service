package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Elenetta17/iris-web-service/internal/httpapi"
)

// TestServerRoutes tests that routes are properly registered
func TestServerRoutes(t *testing.T) {
	// Create a new ServeMux to simulate main's setup
	mux := http.NewServeMux()
	mux.HandleFunc("/", httpapi.FormPage)
	mux.HandleFunc("/hello", httpapi.HelloHandler)

	tests := []struct {
		name string
		path string
	}{
		{
			name: "root path registered",
			path: "/",
		},
		{
			name: "hello path registered",
			path: "/hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.path, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			mux.ServeHTTP(rr, req)

			// Just verify the route is registered (not 404)
			// The actual handler behavior is tested elsewhere
			if rr.Code == http.StatusNotFound {
				t.Errorf("route %s not registered", tt.path)
			}
		})
	}
}
