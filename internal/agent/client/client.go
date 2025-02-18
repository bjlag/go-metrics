package client

import "github.com/bjlag/go-metrics/internal/agent/collector"

type Client interface {
	Send(metrics []*collector.Metric) error
}
