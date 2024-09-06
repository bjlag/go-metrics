package sender_test

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/agent/sender"
	"github.com/bjlag/go-metrics/internal/agent/sender/mock"
)

func TestMetricSender_Send(t *testing.T) {
	type fields struct {
		client func(ctrl *gomock.Controller) *mock.MockClient
	}

	type args struct {
		metric *collector.Metric
	}

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "success",
			fields: fields{
				client: func(ctrl *gomock.Controller) *mock.MockClient {
					request, _ := http.NewRequest(http.MethodPost, "http://127.0.0.1:8080/update/kind/name/1", nil)
					request.Header.Set("Content-Type", "text/plain")

					mockClient := mock.NewMockClient(ctrl)
					mockClient.EXPECT().Do(request).Times(1).Return(&http.Response{StatusCode: http.StatusOK}, nil)

					return mockClient
				},
			},
			args: args{
				metric: collector.NewMetric("kind", "name", 1),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			h := sender.NewHTTPSender(tt.fields.client(ctrl))

			got, _ := h.Send(tt.args.metric)
			defer func() {
				if got.Body != nil {
					_ = got.Body.Close()
				}
			}()

			assert.Equal(t, got.StatusCode, http.StatusOK)
		})
	}
}
