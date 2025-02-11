package rpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/generated/rpc"
)

type MetricSender struct {
	client rpc.MetricServiceClient
}

// NewSender создает RPC клиент.
func NewSender(client *grpc.ClientConn) *MetricSender {
	return &MetricSender{
		client: rpc.NewMetricServiceClient(client),
	}
}

func (s MetricSender) Send(metrics []*collector.Metric) error {
	inMetrics := make([]*rpc.Metric, len(metrics))

	for _, m := range metrics {
		if m == nil {
			continue
		}

		inMetric := &rpc.Metric{
			Id:   m.Name(),
			Type: m.Kind(),
		}

		switch m.Kind() {
		case collector.Gauge:
			value, err := m.GaugeValue()
			if err != nil {
				return err
			}
			inMetric.Value = &value
		case collector.Counter:
			value, err := m.CounterValue()
			if err != nil {
				return err
			}
			inMetric.Delta = &value
		default:
			continue
		}

		inMetrics = append(inMetrics, inMetric)
	}

	out, err := s.client.Updates(context.Background(), &rpc.UpdatesIn{Metrics: inMetrics})
	if err != nil {
	}

	_ = out

	return nil
}
