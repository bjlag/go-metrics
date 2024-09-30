package gauge

import (
	"net/http"
	"strconv"
)

type Handler struct {
	storage storage
	backup  backup
	log     log
}

func NewHandler(storage storage, backup backup, logger log) *Handler {
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
		h.log.Info("metric name not specified")
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	value, err := strconv.ParseFloat(valueMetric, 64)
	if err != nil {
		h.log.WithField("error", err.Error()).Error("invalid metric value")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.storage.SetGauge(nameMetric, value)

	err = h.backup.Create()
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("failed to backup data")
	}
}
