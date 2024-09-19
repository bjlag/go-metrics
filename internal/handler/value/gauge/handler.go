package counter

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	metricNotFoundMsgErr = "Gauge metric '%s' not found"
	writeBodyMsgErr      = "Error while writing body"
)

type Handler struct {
	storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{
		storage: storage,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")

	storeValue, err := h.storage.GetGauge(name)
	if err != nil {
		http.Error(w, fmt.Sprintf(metricNotFoundMsgErr, name), http.StatusNotFound)
		return
	}

	_, err = w.Write([]byte(strconv.FormatFloat(storeValue, 'f', -1, 64)))
	if err != nil {
		http.Error(w, writeBodyMsgErr, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
