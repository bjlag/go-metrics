package rpc

import (
	"context"
	"time"

	"google.golang.org/grpc"

	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/generated/rpc"
	"github.com/bjlag/go-metrics/internal/logger"
)

const timeout = 200 * time.Millisecond

type MetricSender struct {
	client rpc.MetricServiceClient
	log    logger.Logger
}

// NewSender создает RPC клиент.
func NewSender(client *grpc.ClientConn, log logger.Logger) *MetricSender {
	return &MetricSender{
		client: rpc.NewMetricServiceClient(client),
		log:    log,
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := s.client.Updates(ctx, &rpc.UpdatesIn{Metrics: inMetrics})
	if err != nil {
		return err
	}

	s.log.
		WithField("uri", "updates").
		WithField("status", "success").
		Info("Sent RPC request")

	return nil
}
