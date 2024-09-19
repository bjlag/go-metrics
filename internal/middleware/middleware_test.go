package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bjlag/go-metrics/internal/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAllowPostMethodMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		wantStatusCode int
	}{
		{
			name:           "POST",
			method:         http.MethodPost,
			wantStatusCode: http.StatusOK,
		},
		{
			name:           "GET",
			method:         http.MethodGet,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "PUT",
			method:         http.MethodPut,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "PATCH",
			method:         http.MethodPatch,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
		{
			name:           "DELETE",
			method:         http.MethodDelete,
			wantStatusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, "/", nil)

			h := middleware.AllowPostMethodMiddleware(http.HandlerFunc(handler))
			h.ServeHTTP(w, request)

			response := w.Result()
			defer func() {
				_ = response.Body.Close()
			}()

			assert.Equal(t, tt.wantStatusCode, response.StatusCode)
		})

	}
}

func TestFinishRequestMiddleware(t *testing.T) {
	w := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/", nil)

	h := middleware.FinishRequestMiddleware(http.HandlerFunc(handler))
	h.ServeHTTP(w, request)

	response := w.Result()
	defer func() {
		_ = response.Body.Close()
	}()

	assert.Equal(t, "text/plain; charset=utf-8", response.Header.Get("Content-Type"))
}

func handler(_ http.ResponseWriter, _ *http.Request) {}
