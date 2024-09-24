package gauge

import (
	"net/http"
	"strconv"
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
	nameMetric := r.PathValue("name")
	valueMetric := r.PathValue("value")

	if nameMetric == "" {
		h.log.Info("Metric name not specified", nil)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	value, err := strconv.ParseFloat(valueMetric, 64)
	if err != nil {
		h.log.Error("Invalid metric value", nil)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.storage.SetGauge(nameMetric, value)
}
