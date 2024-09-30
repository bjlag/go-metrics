package gauge

import (
	"net/http"
	"strconv"
)

type Handler struct {
	repo   repo
	backup backup
	log    log
}

func NewHandler(repo repo, backup backup, log log) *Handler {
	return &Handler{
		repo:   repo,
		backup: backup,
		log:    log,
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

	h.repo.SetGauge(nameMetric, value)

	err = h.backup.Create()
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("failed to backup data")
	}
}
