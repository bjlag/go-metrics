package update

import (
	"net/http"

	"github.com/bjlag/go-metrics/internal/storage"
)

const (
	invalidMetricTypeMsgErr = "Invalid metric type"
)

type Handler struct {
	storage storage.Interface
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	http.Error(w, invalidMetricTypeMsgErr, http.StatusBadRequest)
}
