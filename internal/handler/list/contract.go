//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package list

import (
	"context"
	"io"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/storage"
)

type renderer interface {
	Render(w io.Writer, name string, data interface{}) error
}

type repo interface {
	GetAllGauges(ctx context.Context) storage.Gauges
	GetAllCounters(ctx context.Context) storage.Counters
}

type log interface {
	WithField(key string, value interface{}) logger.Logger
	Error(msg string)
}
