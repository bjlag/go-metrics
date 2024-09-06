package update

import (
	"net/http"
)

const (
	invalidMetricKindMsgErr = "Invalid metric type"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	kind := r.PathValue("kind")
	if kind == "counter" || kind == "gauge" {
		http.NotFound(w, r)
		return
	}
	http.Error(w, invalidMetricKindMsgErr, http.StatusBadRequest)
}
