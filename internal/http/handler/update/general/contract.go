//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package general

import (
	"context"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/storage"
)

type repo interface {
	SetGauge(ctx context.Context, name string, value float64)
	AddCounter(ctx context.Context, name string, value int64)
	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
	GetAllGauges(ctx context.Context) storage.Gauges
	GetAllCounters(ctx context.Context) storage.Counters
}

type backup interface {
	Create(ctx context.Context) error
}

type log interface {
	WithField(key string, value interface{}) logger.Logger
	Error(msg string)
	Info(msg string)
}
