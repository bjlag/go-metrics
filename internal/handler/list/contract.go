//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package list

import (
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
	Error(msg string, fields map[string]interface{})
}
