package counter

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bjlag/go-metrics/internal/storage/memory"
)

type Handler struct {
	storage Storage
	log     Logger
}

func NewHandler(storage Storage, logger Logger) *Handler {
	return &Handler{
		storage: storage,
		log:     logger,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	storeValue, err := h.storage.GetGauge(name)
	if err != nil {
		var metricNotFoundError *memory.MetricNotFoundError
		if errors.As(err, &metricNotFoundError) {
			h.log.WithField("name", name).
				Info("Gauge metric not found")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		h.log.WithField("error", err.Error()).
			Error("Failed to get gauge value")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(strconv.FormatFloat(storeValue, 'f', -1, 64)))
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("Failed to write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
