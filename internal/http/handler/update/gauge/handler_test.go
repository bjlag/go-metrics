package gauge_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/http/handler/update/gauge"
	"github.com/bjlag/go-metrics/internal/http/handler/update/gauge/mock"
)

func TestHandler_Handle(t *testing.T) {
	type want struct {
		statusCode int
	}

	type fields struct {
		name  string
		value string
	}

	tests := []struct {
		name    string
		storage func(ctrl *gomock.Controller) *mock.Mockrepo
		backup  func(ctrl *gomock.Controller) *mock.Mockbackup
		log     func(ctrl *gomock.Controller) *mock.MockLogger
		fields  fields
		want    want
	}{
		{
			name: "success value is float",
			storage: func(ctrl *gomock.Controller) *mock.Mockrepo {
				mockStorage := mock.NewMockrepo(ctrl)
				mockStorage.EXPECT().SetGauge(gomock.Any(), "test", 1.1).Times(1)

				return mockStorage
			},
			backup: func(ctrl *gomock.Controller) *mock.Mockbackup {
				mockBackup := mock.NewMockbackup(ctrl)
				mockBackup.EXPECT().Create(context.Background()).Times(1)
				return mockBackup
			},
			log: func(ctrl *gomock.Controller) *mock.MockLogger {
				mockLog := mock.NewMockLogger(ctrl)
				mockLog.EXPECT().Info(gomock.Any()).AnyTimes()
				return mockLog
			},
			fields: fields{
				name:  "test",
				value: "1.1",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "success value is int",
			storage: func(ctrl *gomock.Controller) *mock.Mockrepo {
				mockStorage := mock.NewMockrepo(ctrl)
				mockStorage.EXPECT().SetGauge(gomock.Any(), "test", float64(1)).Times(1)

				return mockStorage
			},
			backup: func(ctrl *gomock.Controller) *mock.Mockbackup {
				mockBackup := mock.NewMockbackup(ctrl)
				mockBackup.EXPECT().Create(context.Background()).Times(1)
				return mockBackup
			},
			log: func(ctrl *gomock.Controller) *mock.MockLogger {
				mockLog := mock.NewMockLogger(ctrl)
				mockLog.EXPECT().Info(gomock.Any()).AnyTimes()
				return mockLog
			},
			fields: fields{
				name:  "test",
				value: "1",
			},
			want: want{
				statusCode: http.StatusOK,
			},
		},
		{
			name: "error empty name",
			storage: func(ctrl *gomock.Controller) *mock.Mockrepo {
				mockStorage := mock.NewMockrepo(ctrl)
				mockStorage.EXPECT().SetGauge(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				return mockStorage
			},
			backup: func(ctrl *gomock.Controller) *mock.Mockbackup {
				mockBackup := mock.NewMockbackup(ctrl)
				mockBackup.EXPECT().Create(context.Background()).Times(0)
				return mockBackup
			},
			log: func(ctrl *gomock.Controller) *mock.MockLogger {
				mockLog := mock.NewMockLogger(ctrl)
				mockLog.EXPECT().Info(gomock.Any()).AnyTimes()
				return mockLog
			},
			fields: fields{
				name:  "",
				value: "1",
			},
			want: want{
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "error invalid value is string",
			storage: func(ctrl *gomock.Controller) *mock.Mockrepo {
				mockStorage := mock.NewMockrepo(ctrl)
				mockStorage.EXPECT().SetGauge(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)

				return mockStorage
			},
			backup: func(ctrl *gomock.Controller) *mock.Mockbackup {
				mockBackup := mock.NewMockbackup(ctrl)
				mockBackup.EXPECT().Create(context.Background()).Times(0)
				return mockBackup
			},
			log: func(ctrl *gomock.Controller) *mock.MockLogger {
				mockLog := mock.NewMockLogger(ctrl)
				mockLog.EXPECT().WithField(gomock.Any(), gomock.Any()).Return(mockLog)
				mockLog.EXPECT().Error(gomock.Any()).AnyTimes()
				return mockLog
			},
			fields: fields{
				name:  "test",
				value: "one",
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
			request.SetPathValue("name", tt.fields.name)
			request.SetPathValue("value", tt.fields.value)

			h := http.HandlerFunc(gauge.NewHandler(tt.storage(ctrl), tt.backup(ctrl), tt.log(ctrl)).Handle)
			h.ServeHTTP(w, request)

			response := w.Result()
			defer func() {
				_ = response.Body.Close()
			}()

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}
