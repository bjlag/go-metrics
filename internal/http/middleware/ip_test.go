package middleware_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/http/middleware"
	"github.com/bjlag/go-metrics/internal/mock"
)

func TestCheckRealIPMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		logger        func(ctrl *gomock.Controller) *mock.MockLogger
		realIP        string
		trustedSubnet string
		wantStatus    int
	}{
		{
			name:          "success",
			logger:        mock.NewMockLogger,
			realIP:        "192.168.1.1",
			trustedSubnet: "192.168.1.0/24",
			wantStatus:    http.StatusOK,
		},
		{
			name:          "trusted subnet is empty",
			logger:        mock.NewMockLogger,
			realIP:        "192.168.1.1",
			trustedSubnet: "",
			wantStatus:    http.StatusOK,
		},
		{
			name: "rejected",
			logger: func(ctrl *gomock.Controller) *mock.MockLogger {
				mockLogger := mock.NewMockLogger(ctrl)
				mockLogger.EXPECT().WithField("IP", "192.168.2.1").Return(mockLogger)
				mockLogger.EXPECT().Error(gomock.Any())
				return mockLogger
			},
			realIP:        "192.168.2.1",
			trustedSubnet: "192.168.1.0/24",
			wantStatus:    http.StatusForbidden,
		},
		{
			name: "rejected, X-Real-IP is empty",
			logger: func(ctrl *gomock.Controller) *mock.MockLogger {
				mockLogger := mock.NewMockLogger(ctrl)
				mockLogger.EXPECT().Error(gomock.Any())
				return mockLogger
			},
			realIP:        "",
			trustedSubnet: "192.168.1.0/24",
			wantStatus:    http.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			w := httptest.NewRecorder()
			request := httptest.NewRequest("POST", "/url", nil)
			request.Header.Set("X-Real-IP", tt.realIP)

			var (
				err   error
				ipNet *net.IPNet
			)

			if tt.trustedSubnet != "" {
				_, ipNet, err = net.ParseCIDR(tt.trustedSubnet)
				assert.NoError(t, err)
			}

			h := middleware.CheckRealIPMiddleware(ipNet, tt.logger(ctrl))(http.HandlerFunc(handlerHeaderResponse))
			h.ServeHTTP(w, request)

			assert.Equal(t, w.Code, tt.wantStatus)
		})
	}
}
