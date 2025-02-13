//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package batch

import (
	"context"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/storage"
)

type repo interface {
	SetGauges(ctx context.Context, gauges []storage.Gauge) error
	AddCounters(ctx context.Context, counters []storage.Counter) error
}

type backup interface {
	Create(ctx context.Context) error
}

type log interface {
	WithError(err error) logger.Logger
	Error(msg string)
	Info(msg string)
}
