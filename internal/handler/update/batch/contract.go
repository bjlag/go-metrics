//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package batch

import (
	"context"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/storage"
)

type repo interface {
	SetGauge(ctx context.Context, name string, value float64)
	SetGauges(ctx context.Context, gauges []storage.Gauge) error
	AddCounter(ctx context.Context, name string, value int64)
	AddCounters(ctx context.Context, counters []storage.Counter) error
	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
	GetAllGauges(ctx context.Context) storage.Gauges
	GetAllCounters(ctx context.Context) storage.Counters
}

type backup interface {
	Create(ctx context.Context) error
}

type log interface {
	WithError(err error) logger.Logger
	Error(msg string)
	Info(msg string)
}
