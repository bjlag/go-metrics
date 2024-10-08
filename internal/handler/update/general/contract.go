//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package general

import (
	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/storage"
)

type repo interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetAllGauges() storage.Gauges
	GetAllCounters() storage.Counters
}

type backup interface {
	Create() error
}

type log interface {
	WithField(key string, value interface{}) logger.Logger
	Error(msg string)
	Info(msg string)
}
