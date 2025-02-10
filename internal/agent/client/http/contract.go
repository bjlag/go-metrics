//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go
//go:generate mockgen -package mock -destination mock/logger_mock.go github.com/bjlag/go-metrics/internal/logger Logger

package http

import "github.com/bjlag/go-metrics/internal/logger"

type log interface {
	WithField(key string, value interface{}) logger.Logger
	Error(msg string)
}
