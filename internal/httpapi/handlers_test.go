package httpapi

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// errReadCloser is a small helper to simulate a read error from the request body.
type errReadCloser struct{}

func (errReadCloser) Read(p []byte) (int, error) { return 0, errors.New("read error") }
func (errReadCloser) Close() error               { return nil }

func TestHelloHandler(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		form         url.Values      // form values to send (nil -> no body)
		body         io.ReadCloser   // optional custom body (overrides form)
		contentType  string
		wantCode     int
		wantBody     string
		wantCTPrefix string // optional: check response Content-Type starts with this
	}{
		{
			name:     "wrong method",
			method:   http.MethodGet,
			wantCode: http.StatusMethodNotAllowed,
		},
		{
			name:     "default name",
			method:   http.MethodPost,
			form:     url.Values{},
			contentType: "application/x-www-form-urlencoded",
			wantCode: http.StatusOK,
			wantBody: "Hello World!",
			wantCTPrefix: "text/plain",
		},
		{
			name:     "with name",
			method:   http.MethodPost,
			form:     url.Values{"name": {"Alice"}},
			contentType: "application/x-www-form-urlencoded",
			wantCode: http.StatusOK,
			wantBody: "Hello Alice!",
			wantCTPrefix: "text/plain",
		},
		{
			name:    "parse form error",
			method:  http.MethodPost,
			// supply a body that returns error when read to force r.ParseForm() to fail
			body:        errReadCloser{},
			contentType: "application/x-www-form-urlencoded",
			wantCode:    http.StatusBadRequest,
			wantBody:    "Invalid form\n", // http.Error appends newline
		},
	}

	for _, tt := range tests {
		tt := tt // capture range variable for running tests in parallel if desired
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel() // uncomment if tests are safe to run in parallel

			var reqBody io.ReadCloser
			if tt.body != nil {
				reqBody = tt.body
			} else if tt.form != nil {
				reqBody = io.NopCloser(strings.NewReader(tt.form.Encode()))
			}

			req := httptest.NewRequest(tt.method, "/hello", reqBody)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			// Use ResponseRecorder to capture the handler's response
			rr := httptest.NewRecorder()
			HelloHandler(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.wantCode {
				t.Fatalf("status: want %d, got %d", tt.wantCode, res.StatusCode)
			}

			if tt.wantCTPrefix != "" {
				ct := res.Header.Get("Content-Type")
				if !strings.HasPrefix(ct, tt.wantCTPrefix) {
					t.Fatalf("Content-Type: want prefix %q, got %q", tt.wantCTPrefix, ct)
				}
			}

			// Read body if we expect one (empty wantBody means we don't check)
			if tt.wantBody != "" {
				b, err := io.ReadAll(res.Body)
				if err != nil {
					t.Fatalf("reading body: %v", err)
				}
				got := string(b)
				if got != tt.wantBody {
					t.Fatalf("body: want %q, got %q", tt.wantBody, got)
				}
			}
		})
	}
}

