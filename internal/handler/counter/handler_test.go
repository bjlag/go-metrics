package counter_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/handler/counter"
	"github.com/bjlag/go-metrics/internal/handler/counter/mock"
)

func Test_Handle(t *testing.T) {
	type want struct {
		statusCode int
	}

	type fields struct {
		name  string
		value string
	}

	tests := []struct {
		name    string
		storage func(ctrl *gomock.Controller) *mock.MockStorage
		fields  fields
		want    want
	}{
		{
			name: "success",
			storage: func(ctrl *gomock.Controller) *mock.MockStorage {
				mockStorage := mock.NewMockStorage(ctrl)
				mockStorage.EXPECT().AddCounter("test", int64(1)).Times(1)

				return mockStorage
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
			storage: func(ctrl *gomock.Controller) *mock.MockStorage {
				mockStorage := mock.NewMockStorage(ctrl)
				mockStorage.EXPECT().AddCounter(gomock.Any(), gomock.Any()).Times(0)

				return mockStorage
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
			name: "error invalid value is float",
			storage: func(ctrl *gomock.Controller) *mock.MockStorage {
				mockStorage := mock.NewMockStorage(ctrl)
				mockStorage.EXPECT().AddCounter(gomock.Any(), gomock.Any()).Times(0)

				return mockStorage
			},
			fields: fields{
				name:  "test",
				value: "1.1",
			},
			want: want{
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "error invalid value is string",
			storage: func(ctrl *gomock.Controller) *mock.MockStorage {
				mockStorage := mock.NewMockStorage(ctrl)
				mockStorage.EXPECT().AddCounter(gomock.Any(), gomock.Any()).Times(0)

				return mockStorage
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

			h := http.HandlerFunc(counter.NewHandler(tt.storage(ctrl)).Handle)
			h.ServeHTTP(w, request)

			response := w.Result()
			defer func() {
				_ = response.Body.Close()
			}()

			assert.Equal(t, tt.want.statusCode, response.StatusCode)
		})
	}
}
