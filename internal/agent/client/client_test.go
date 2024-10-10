package client_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/agent/client"
	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/model"
)

func TestMetricSender_Send(t *testing.T) {
	type args struct {
		metric []*collector.Metric
	}

	tests := []struct {
		name       string
		args       args
		server     *httptest.Server
		wantStatus int
		wantErr    bool
	}{
		{
			name: "success counter",
			args: args{
				metric: []*collector.Metric{collector.NewMetric("counter", "name", 1)},
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/update/", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "POST", r.Method)

				var err error
				var buf bytes.Buffer

				_, err = buf.ReadFrom(r.Body)
				assert.NoError(t, err)

				in := &model.UpdateIn{}
				err = json.Unmarshal(buf.Bytes(), in)
				assert.NoError(t, err)

				assert.Equal(t, "name", in.ID)
				assert.Equal(t, "counter", in.MType)
				assert.Equal(t, int64(1), *in.Delta)
				assert.Nil(t, in.Value)
			})),
			wantStatus: http.StatusOK,
		},
		{
			name: "success gauge",
			args: args{
				metric: []*collector.Metric{collector.NewMetric("gauge", "name", 1.1)},
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/update/", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "POST", r.Method)

				var err error
				var buf bytes.Buffer

				_, err = buf.ReadFrom(r.Body)
				assert.NoError(t, err)

				in := &model.UpdateIn{}
				err = json.Unmarshal(buf.Bytes(), in)
				assert.NoError(t, err)

				assert.Equal(t, "name", in.ID)
				assert.Equal(t, "gauge", in.MType)
				assert.Equal(t, 1.1, *in.Value)
				assert.Nil(t, in.Delta)
			})),
			wantStatus: http.StatusOK,
		},
		{
			name: "error unknown metric type",
			args: args{
				metric: []*collector.Metric{collector.NewMetric("unknown_type", "name", 1.1)},
			},
			server: httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "/update/", r.URL.Path)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "POST", r.Method)

				var err error
				var buf bytes.Buffer

				_, err = buf.ReadFrom(r.Body)
				assert.NoError(t, err)

				in := &model.UpdateIn{}
				err = json.Unmarshal(buf.Bytes(), in)
				assert.NoError(t, err)

				assert.Equal(t, "name", in.ID)
				assert.Equal(t, "gauge", in.MType)
				assert.Equal(t, 1.1, *in.Value)
				assert.Nil(t, in.Delta)
			})),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Skip()

			defer tt.server.Close()

			parts := strings.Split(tt.server.Listener.Addr().String(), ":")
			port, _ := strconv.Atoi(parts[1])

			c := client.NewHTTPSender(parts[0], port)

			got, err := c.Send(tt.args.metric)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.Equal(t, tt.wantStatus, got.StatusCode())
		})
	}
}
