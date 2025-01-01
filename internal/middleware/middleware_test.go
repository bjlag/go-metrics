package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/middleware"
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

func handler(_ http.ResponseWriter, _ *http.Request) {}
