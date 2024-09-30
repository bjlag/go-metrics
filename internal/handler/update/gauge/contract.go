//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go
//go:generate mockgen -package mock -destination mock/logger_mock.go github.com/bjlag/go-metrics/internal/logger Logger

package gauge

import internalLogger "github.com/bjlag/go-metrics/internal/logger"

type storage interface {
	SetGauge(name string, value float64)
}

type backup interface {
	Create() error
}

type log interface {
	WithField(key string, value interface{}) internalLogger.Logger
	Error(msg string)
	Info(msg string)
}
