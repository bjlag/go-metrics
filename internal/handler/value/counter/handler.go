package counter

import (
	"errors"
	"fmt"
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

	storedValue, err := h.storage.GetCounter(name)
	if err != nil {
		var metricNotFoundError *memory.MetricNotFoundError
		if errors.As(err, &metricNotFoundError) {
			h.log.Info(fmt.Sprintf("Counter metric not found: %s", name), nil)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		h.log.Error(fmt.Sprintf("Failed to get counter value: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte(strconv.FormatInt(storedValue, 10)))
	if err != nil {
		h.log.Error(fmt.Sprintf("Failed to write response: %s", err.Error()), nil)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
