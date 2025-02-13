package rpc

import (
	"context"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	clientIP "github.com/bjlag/go-metrics/internal/agent/client"
	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/generated/rpc"
	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/rpc/interceptor"
	"github.com/bjlag/go-metrics/internal/securety/signature"
)

const timeout = 200 * time.Millisecond

type MetricSender struct {
	conn     *grpc.ClientConn
	client   rpc.MetricServiceClient
	clientIP net.IP
	sign     *signature.SignManager
}

// NewSender создает RPC клиент.
func NewSender(addr string, sign *signature.SignManager, log logger.Logger) *MetricSender {
	conn, err := grpc.NewClient(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			interceptor.LoggerClientInterceptor(log),
			interceptor.RealIPInterceptor,
		),
	)
	if err != nil {
		panic(err)
	}

	return &MetricSender{
		conn:     conn,
		client:   rpc.NewMetricServiceClient(conn),
		clientIP: clientIP.GetOutboundIP(),
		sign:     sign,
	}
}

func (s MetricSender) Close() error {
	return s.conn.Close()
}

func (s MetricSender) Send(metrics []*collector.Metric) error {
	inMetrics := make([]*rpc.Metric, 0, len(metrics))

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

	//mdMap := map[string]string{
	//	"RealIP": s.clientIP.String(),
	//}

	//if s.sign.Enable() {
	//	jsonb, err := json.Marshal(inMetrics)
	//	if err != nil {
	//		return fmt.Errorf("failed to marshal metric: %s", err)
	//	}
	//
	//	mdMap["HashSHA256"] = s.sign.Sing(jsonb)
	//}

	//ctx = metadata.NewOutgoingContext(ctx, metadata.New(mdMap))

	_, err := s.client.Updates(ctx, &rpc.UpdatesIn{Metrics: inMetrics})
	if err != nil {
		return err
	}

	return nil
}
