package update_counter

import (
	"log"
	"net/http"
	"strconv"

	"github.com/bjlag/go-metrics/internal/storage"
)

const (
	emptyNameMetricMsgErr    = "Metric name not specified"
	invalidMetricValueMsgErr = "Invalid metric value"
)

type Handler struct {
	storage storage.Interface
}

func NewHandler(storage storage.Interface) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	nameMetric := r.PathValue("name")
	valueMetric := r.PathValue("value")

	if nameMetric == "" {
		http.Error(w, emptyNameMetricMsgErr, http.StatusNotFound)
		return
	}

	value, err := strconv.ParseInt(valueMetric, 10, 64)
	if err != nil {
		http.Error(w, invalidMetricValueMsgErr, http.StatusBadRequest)
		return
	}

	h.storage.AddCounter(nameMetric, value)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	log.Printf("Metric received: '%s', name '%s', value '%s'\n", "counter", nameMetric, valueMetric)
}
