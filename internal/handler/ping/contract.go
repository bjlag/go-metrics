package ping

import "github.com/bjlag/go-metrics/internal/logger"

type db interface {
	Ping() error
}

type log interface {
	WithError(err error) logger.Logger
	Error(msg string)
}
