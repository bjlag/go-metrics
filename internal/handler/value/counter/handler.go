package counter

import (
	"fmt"
	"net/http"
	"strconv"
)

const (
	metricNotFoundMsgErr = "Counter metric '%s' not found"
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

	storedValue, err := h.storage.GetCounter(name)
	if err != nil {
		http.Error(w, fmt.Sprintf(metricNotFoundMsgErr, name), http.StatusNotFound)
		return
	}

	_, err = w.Write([]byte(strconv.FormatInt(storedValue, 10)))
	if err != nil {
		http.Error(w, writeBodyMsgErr, http.StatusInternalServerError)
	}
}
