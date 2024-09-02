package update_counter

import (
	"net/http"
	"strconv"

	"github.com/bjlag/go-metrics/internal/storage"
)

const (
	invalidTypeValueMsgErr = "Invalid type value of metric"
)

type Handler struct {
	storage storage.Interface
}

func NewHandler(storage storage.Interface) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request, nameMetric, valueMetric string) {
	value, err := strconv.ParseInt(valueMetric, 10, 64)
	if err != nil {
		http.Error(w, invalidTypeValueMsgErr, http.StatusBadRequest)
		return
	}

	h.storage.AddCounter(nameMetric, value)
}
