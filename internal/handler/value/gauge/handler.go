package counter

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bjlag/go-metrics/internal/storage"
)

// Handler обработчик HTTP запроса на получение значения метрики типа Gauge.
type Handler struct {
	repo repo
	log  log
}

// NewHandler создает обработчик.
func NewHandler(repo repo, log log) *Handler {
	return &Handler{
		repo: repo,
		log:  log,
	}
}

// Handle обрабатывает HTTP запрос.
func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	storeValue, err := h.repo.GetGauge(r.Context(), name)
	if err != nil {
		var metricNotFoundError *storage.NotFoundError
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
