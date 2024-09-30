//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package counter

import "github.com/bjlag/go-metrics/internal/logger"

type Storage interface {
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}

type Logger interface {
	WithField(key string, value interface{}) logger.Logger
	Error(msg string)
	Info(msg string)
}
