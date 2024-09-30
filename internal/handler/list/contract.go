//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package list

import (
	"github.com/bjlag/go-metrics/internal/logger"
	"io"

	"github.com/bjlag/go-metrics/internal/storage"
)

type Renderer interface {
	Render(w io.Writer, name string, data interface{}) error
}

type Storage interface {
	GetAllGauges() storage.Gauges
	GetAllCounters() storage.Counters
}

type Logger interface {
	WithField(key string, value interface{}) logger.Logger
	Error(msg string)
}
