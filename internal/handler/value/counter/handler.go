package counter

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bjlag/go-metrics/internal/storage/memory"
)

type Handler struct {
	repo repo
	log  log
}

func NewHandler(repo repo, log log) *Handler {
	return &Handler{
		repo: repo,
		log:  log,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	storedValue, err := h.repo.GetCounter(r.Context(), name)
	if err != nil {
		var metricNotFoundError *memory.MetricNotFoundError
		if errors.As(err, &metricNotFoundError) {
			h.log.WithField("name", name).
				Info("Counter metric not found")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		h.log.WithField("error", err.Error()).
			Error("Failed to get counter value")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(strconv.FormatInt(storedValue, 10)))
	if err != nil {
		h.log.WithField("error", err.Error()).
			Error("Failed to write response")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
