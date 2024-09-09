package client_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bjlag/go-metrics/internal/agent/client"
	"github.com/bjlag/go-metrics/internal/agent/collector"
)

func TestMetricSender_Send(t *testing.T) {
	type args struct {
		metric *collector.Metric
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "success",
			args: args{
				metric: collector.NewMetric("kind", "name", 1),
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/update/kind/name/1", r.URL.Path)
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
	}))
	defer server.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts := strings.Split(server.Listener.Addr().String(), ":")
			port, _ := strconv.Atoi(parts[1])

			c := client.NewHTTPSender(parts[0], port)

			got, _ := c.Send(tt.args.metric)

			assert.Equal(t, got.StatusCode(), http.StatusOK)
		})
	}
}
