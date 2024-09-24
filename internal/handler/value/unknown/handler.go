package unknown

import (
	"fmt"
	"net/http"
)

type Handler struct {
	log Logger
}

func NewHandler(logger Logger) *Handler {
	return &Handler{
		log: logger,
	}
}

func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	h.log.Info(fmt.Sprintf("Invalid metric type %s, url %s", r.PathValue("kind"), r.URL.Path), nil)
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
