package counter

import (
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
		h.log.Info("Metric name not specified")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	value, err := strconv.ParseInt(valueMetric, 10, 64)
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("Invalid metric value")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.storage.AddCounter(nameMetric, value)

	err = h.backup.Create()
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("Failed to backup data")
	}
}
