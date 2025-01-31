package client_test

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/agent/client"
	"github.com/bjlag/go-metrics/internal/agent/client/mock"
	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/agent/limiter"
	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/securety/crypt"
	"github.com/bjlag/go-metrics/internal/securety/signature"
)

func TestMetricSender_Send(t *testing.T) {
	type args struct {
		metric []*collector.Metric
	}

	tests := []struct {
		name       string
		args       args
		server     *httptest.Server
		log        func(ctrl *gomock.Controller) *mock.MockLogger
		wantStatus int
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				metric: []*collector.Metric{
					collector.NewMetric("counter", "counter_name", 1),
					collector.NewMetric("gauge", "gauge_name", 2),
				},
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/updates/", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "POST", r.Method)

				body := decompress(r.Body)

				var in []model.UpdateIn
				err := json.Unmarshal(body, &in)
				assert.NoError(t, err)

				assert.Len(t, in, 2)

				assert.Equal(t, "counter", in[0].MType)
				assert.Equal(t, "counter_name", in[0].ID)
				assert.Equal(t, int64(1), *in[0].Delta)
				assert.Nil(t, in[0].Value)

				assert.Equal(t, "gauge", in[1].MType)
				assert.Equal(t, "gauge_name", in[1].ID)
				assert.Equal(t, float64(2), *in[1].Value)
				assert.Nil(t, in[1].Delta)
			})),
			log: func(ctrl *gomock.Controller) *mock.MockLogger {
				mockLog := mock.NewMockLogger(ctrl)
				mockLog.EXPECT().WithField(gomock.Any(), gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any())
				return mockLog
			},
			wantStatus: http.StatusOK,
		},
		{
			name: "error unknown metric type",
			args: args{
				metric: []*collector.Metric{collector.NewMetric("unknown_type", "name", 1.1)},
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/updates/", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "POST", r.Method)

				body := decompress(r.Body)

				var in []model.UpdateIn
				err := json.Unmarshal(body, &in)
				assert.NoError(t, err)

				assert.Len(t, in, 0)
			})),
			log: func(ctrl *gomock.Controller) *mock.MockLogger {
				mockLog := mock.NewMockLogger(ctrl)
				mockLog.EXPECT().WithField(gomock.Any(), gomock.Any()).Return(mockLog).AnyTimes()
				mockLog.EXPECT().Info(gomock.Any())
				return mockLog
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			defer tt.server.Close()

			parts := strings.Split(tt.server.Listener.Addr().String(), ":")
			port, _ := strconv.Atoi(parts[1])

			encryptManager, _ := crypt.NewEncryptManager("")

			c := client.NewHTTPSender(
				parts[0],
				port,
				signature.NewSignManager("secretKey"),
				encryptManager,
				limiter.NewRateLimiter(1),
				tt.log(ctrl),
			)

			err := c.Send(tt.args.metric)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
		})
	}
}

func decompress(r io.Reader) []byte {
	zr, _ := gzip.NewReader(r)

	var zb bytes.Buffer
	_, _ = zb.ReadFrom(zr)

	return zb.Bytes()
}
