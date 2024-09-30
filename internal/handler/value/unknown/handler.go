package unknown

import (
	"net/http"
)

type Handler struct {
	log log
}

func NewHandler(log log) *Handler {
	return &Handler{
		log: log,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	h.log.WithField("type", r.PathValue("kind")).
		WithField("url", r.URL.Path).
		Info("Invalid metric type")
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
