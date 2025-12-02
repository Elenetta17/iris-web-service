package httpapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func runHelloRequest(t *testing.T, method string, form url.Values, contentType string) *httptest.ResponseRecorder {
	t.Helper()

	var body io.Reader
	if form != nil {
		body = strings.NewReader(form.Encode())
	}

	req, err := http.NewRequest(method, "/hello", body)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(HelloHandler).ServeHTTP(rr, req)
	return rr
}

func TestFormPage(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	http.HandlerFunc(FormPage).ServeHTTP(rr, req)

	if got, want := rr.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}

	if got, want := rr.Header().Get("Content-Type"), "text/html"; got != want {
		t.Fatalf("content-type = %q, want %q", got, want)
	}

	body := rr.Body.String()
	for _, element := range []string{
		`<form`,
		`action="/hello"`,
		`method="POST"`,
		`<input`,
		`name="name"`,
		`<button`,
		`type="submit"`,
	} {
		if !strings.Contains(body, element) {
			t.Errorf("response body missing required element: %q", element)
		}
	}
}

func TestHelloHandler_Success(t *testing.T) {
	form := url.Values{"name": {"Alice"}}
	rr := runHelloRequest(t, http.MethodPost, form, "application/x-www-form-urlencoded")

	if got, want := rr.Code, http.StatusOK; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}

	if got, want := rr.Header().Get("Content-Type"), "text/plain"; got != want {
		t.Fatalf("content-type = %q, want %q", got, want)
	}

	if got, want := rr.Body.String(), "Hello Alice!"; got != want {
		t.Errorf("body = %q, want %q", got, want)
	}
}

func TestHelloHandler_DefaultWorld(t *testing.T) {
	tests := []struct {
		name string
		form url.Values
	}{
		{"empty name", url.Values{"name": {""}}},
		{"no name", url.Values{}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rr := runHelloRequest(t, http.MethodPost, tc.form, "application/x-www-form-urlencoded")

			if got, want := rr.Code, http.StatusOK; got != want {
				t.Fatalf("status = %d, want %d", got, want)
			}

			if got, want := rr.Body.String(), "Hello World!"; got != want {
				t.Errorf("body = %q, want %q", got, want)
			}
		})
	}
}

func TestHelloHandler_MethodNotAllowed(t *testing.T) {
	for _, method := range []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	} {
		t.Run(method, func(t *testing.T) {
			rr := runHelloRequest(t, method, nil, "")

			if got, want := rr.Code, http.StatusMethodNotAllowed; got != want {
				t.Fatalf("status = %d, want %d", got, want)
			}
			if !strings.Contains(rr.Body.String(), "Method not allowed") {
				t.Errorf("body %q does not contain %q", rr.Body.String(), "Method not allowed")
			}
		})
	}
}

func TestHelloHandler_InvalidForm(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/hello", strings.NewReader("invalid"))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=")

	rr := httptest.NewRecorder()
	http.HandlerFunc(HelloHandler).ServeHTTP(rr, req)

	if got, want := rr.Code, http.StatusBadRequest; got != want {
		t.Fatalf("status = %d, want %d", got, want)
	}
	if !strings.Contains(rr.Body.String(), "Invalid form") {
		t.Errorf("body %q does not contain %q", rr.Body.String(), "Invalid form")
	}
}

func TestHelloHandler_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"spaces", "John Doe", "Hello John Doe!"},
		{"unicode", "José", "Hello José!"},
		{"numbers", "User123", "Hello User123!"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			form := url.Values{"name": {tc.input}}
			rr := runHelloRequest(t, http.MethodPost, form, "application/x-www-form-urlencoded")

			if got := rr.Body.String(); got != tc.expected {
				t.Errorf("body = %q, want %q", got, tc.expected)
			}
		})
	}
}

