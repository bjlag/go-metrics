package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/bjlag/go-metrics/internal/http/middleware"
	"github.com/bjlag/go-metrics/internal/mock"
)

func TestLogMiddleware(t *testing.T) {
	tests := []struct {
		name   string
		logger func(ctrl *gomock.Controller) *mock.MockLogger
	}{
		{
			name: "log_request",
			logger: func(ctrl *gomock.Controller) *mock.MockLogger {
				mockLogger := mock.NewMockLogger(ctrl)
				mockLogger.EXPECT().WithField("uri", "/url").Return(mockLogger).AnyTimes()
				mockLogger.EXPECT().WithField("method", "POST").Return(mockLogger).AnyTimes()
				mockLogger.EXPECT().WithField("content_type", "application/json").Return(mockLogger).AnyTimes()
				mockLogger.EXPECT().WithField("duration", gomock.Any()).Return(mockLogger).AnyTimes()
				mockLogger.EXPECT().WithField("status", 200).Return(mockLogger).AnyTimes()
				mockLogger.EXPECT().WithField("size", 7).Return(mockLogger).AnyTimes()
				mockLogger.EXPECT().Info("Got request").AnyTimes()

				return mockLogger
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			w := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/url", nil)
			request.Header.Set("Content-Type", "application/json")

			h := middleware.LogMiddleware(tt.logger(ctrl))(http.HandlerFunc(handlerLogRequest))
			h.ServeHTTP(w, request)
		})
	}
}

func handlerLogRequest(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("handler"))
	w.WriteHeader(http.StatusOK)
}
