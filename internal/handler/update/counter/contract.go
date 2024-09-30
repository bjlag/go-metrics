//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package counter

import "github.com/bjlag/go-metrics/internal/logger"

type Storage interface {
	AddCounter(name string, value int64)
}

type Backup interface {
	Create() error
}

type Logger interface {
	WithField(key string, value interface{}) logger.Logger
	Error(msg string)
	Info(msg string)
}
