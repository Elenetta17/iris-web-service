package metrics

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestMiddleware(t *testing.T) {
	// Reset the counter before testing
	HttpRequestsTotal.Reset()

	// Create a simple test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Wrap it with our middleware
	wrappedHandler := Middleware(testHandler)

	// Make a test request
	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	wrappedHandler.ServeHTTP(rec, req)

	// Verify the counter was incremented
	count := testutil.ToFloat64(HttpRequestsTotal.WithLabelValues("/test", "GET"))
	if count != 1 {
		t.Errorf("expected counter to be 1, got %f", count)
	}

	// Make another request to same endpoint
	wrappedHandler.ServeHTTP(rec, req)
	count = testutil.ToFloat64(HttpRequestsTotal.WithLabelValues("/test", "GET"))
	if count != 2 {
		t.Errorf("expected counter to be 2, got %f", count)
	}

	// Make a request to different endpoint
	req2 := httptest.NewRequest("POST", "/hello", nil)
	wrappedHandler.ServeHTTP(rec, req2)
	count = testutil.ToFloat64(HttpRequestsTotal.WithLabelValues("/hello", "POST"))
	if count != 1 {
		t.Errorf("expected counter to be 1, got %f", count)
	}
}
