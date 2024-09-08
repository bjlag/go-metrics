package client_test

import (
	"net/http"
	"net/http/httptest"
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
			c := client.NewHTTPSender(server.URL)

			got, _ := c.Send(tt.args.metric)

			assert.Equal(t, got.StatusCode(), http.StatusOK)
		})
	}
}
