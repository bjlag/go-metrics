package unknown

import (
	"log"
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
	log.Printf("Invalid metric type %s, url %s", r.PathValue("kind"), r.URL.Path)
	http.Error(w, invalidMetricKindMsgErr, http.StatusBadRequest)
}
