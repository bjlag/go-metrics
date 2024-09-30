//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go
//go:generate mockgen -package mock -destination mock/logger_mock.go github.com/bjlag/go-metrics/internal/logger Logger

package counter

import "github.com/bjlag/go-metrics/internal/logger"

type storage interface {
	AddCounter(name string, value int64)
}

type backup interface {
	Create() error
}

type log interface {
	WithField(key string, value interface{}) logger.Logger
	Error(msg string)
	Info(msg string)
}
