//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package unknown

import "github.com/bjlag/go-metrics/internal/logger"

type Logger interface {
	WithField(key string, value interface{}) logger.Logger
	Info(msg string)
}
