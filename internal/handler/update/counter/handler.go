package counter

import (
	"net/http"
	"strconv"
)

// Handler обработчик HTTP запроса на обновление метрики типа Counter.
type Handler struct {
	repo   repo
	backup backup
	log    log
}

// NewHandler создает обработчик.
func NewHandler(repo repo, backup backup, log log) *Handler {
	return &Handler{
		repo:   repo,
		backup: backup,
		log:    log,
	}
}

// Handle обрабатывает HTTP запрос.
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

	h.repo.AddCounter(r.Context(), nameMetric, value)

	err = h.backup.Create(r.Context())
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("Failed to backup data")
	}
}
