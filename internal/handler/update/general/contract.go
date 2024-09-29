//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package general

import (
	"github.com/bjlag/go-metrics/internal/storage"
)

type Storage interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
	GetAllGauges() storage.Gauges
	GetAllCounters() storage.Counters
}

type Backup interface {
	Create() error
}

type Logger interface {
	Error(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
}
