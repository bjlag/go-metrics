package gauge

import (
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	storage Storage
	backup  Backup
	log     Logger
}

func NewHandler(storage Storage, backup Backup, logger Logger) *Handler {
	return &Handler{
		storage: storage,
		backup:  backup,
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

	err = h.backup.Create()
	if err != nil {
		h.log.Error(fmt.Sprintf("Failed to backup data: %s", err.Error()), nil)
	}
}
