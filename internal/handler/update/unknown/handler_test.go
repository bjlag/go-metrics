package unknown_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/handler/update/unknown"
	"github.com/bjlag/go-metrics/internal/handler/update/unknown/mock"
)

func TestHandler_Handle(t *testing.T) {
	type want struct {
		statusCode int
	}

	type fields struct {
		kind string
		log  func(ctrl *gomock.Controller) *mock.MockLogger
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
				log: func(ctrl *gomock.Controller) *mock.MockLogger {
					mockLog := mock.NewMockLogger(ctrl)
					mockLog.EXPECT().WithField(gomock.Any(), gomock.Any()).Return(mockLog).AnyTimes()
					mockLog.EXPECT().Info(gomock.Any()).AnyTimes()
					return mockLog
				},
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			w := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodPost, "/", nil)
			request.SetPathValue("kind", tt.fields.kind)

			h := http.HandlerFunc(unknown.NewHandler(tt.fields.log(ctrl)).Handle)
			h.ServeHTTP(w, request)

			response := w.Result()
			defer func() {
				_ = response.Body.Close()
			}()

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}
