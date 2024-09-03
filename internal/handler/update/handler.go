package update

import (
	"net/http"
)

const (
	invalidMetricTypeMsgErr = "Invalid metric type"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	typeMetric := r.PathValue("type")
	if typeMetric == "counter" || typeMetric == "gauge" {
		http.NotFound(w, r)
		return
	}
	http.Error(w, invalidMetricTypeMsgErr, http.StatusBadRequest)
}
