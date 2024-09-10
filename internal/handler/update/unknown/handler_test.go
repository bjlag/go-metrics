package unknown_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/handler/update/unknown"
)

func TestHandler_Handle(t *testing.T) {
	type want struct {
		statusCode int
	}

	type fields struct {
		kind string
	}

	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			name: "other some thing",
			fields: fields{
				kind: "other",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodPost, "/", nil)
			request.SetPathValue("kind", tt.fields.kind)

			h := http.HandlerFunc(unknown.NewHandler().Handle)
			h.ServeHTTP(w, request)

			response := w.Result()
			defer func() {
				_ = response.Body.Close()
			}()

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}
