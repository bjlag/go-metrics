//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package general

import (
	"context"

	"github.com/bjlag/go-metrics/internal/logger"
)

type repo interface {
	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
}

type log interface {
	WithField(key string, value interface{}) logger.Logger
	Error(msg string)
	Info(msg string)
}
