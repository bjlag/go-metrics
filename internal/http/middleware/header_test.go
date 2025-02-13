package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/http/middleware"
)

func TestHeaderResponseMiddleware(t *testing.T) {
	type header struct {
		key    string
		values []string
	}
	tests := []struct {
		name   string
		header header
	}{
		{
			name: "one",
			header: header{
				key:    "Content-Type",
				values: []string{"application/json"},
			},
		},
		{
			name: "several",
			header: header{
				key:    "Content-Type",
				values: []string{"application/json", "utf-8"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/url", nil)

			h := middleware.HeaderResponseMiddleware(tt.header.key, tt.header.values...)(http.HandlerFunc(handlerHeaderResponse))
			h.ServeHTTP(w, request)

			assert.Equal(t, strings.Join(tt.header.values, "; "), w.Header().Get(tt.header.key))
		})
	}
}

func handlerHeaderResponse(w http.ResponseWriter, _ *http.Request) {}
