package counter

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bjlag/go-metrics/internal/storage/memory"
)

const (
	metricNotFoundMsgErr = "Gauge metric '%s' not found"
	writeBodyMsgErr      = "Error while writing body"
)

type Handler struct {
	storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	storeValue, err := h.storage.GetGauge(name)
	if err != nil {
		var metricNotFoundError *memory.MetricNotFoundError
		if errors.As(err, &metricNotFoundError) {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(strconv.FormatFloat(storeValue, 'f', -1, 64)))
	if err != nil {
		http.Error(w, writeBodyMsgErr, http.StatusInternalServerError)
	}
}
