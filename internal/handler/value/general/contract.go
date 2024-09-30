//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package general

import "github.com/bjlag/go-metrics/internal/logger"

type repo interface {
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}

type log interface {
	WithField(key string, value interface{}) logger.Logger
	Error(msg string)
	Info(msg string)
}
