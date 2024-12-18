package ping

import (
	"net/http"
)

type Handler struct {
	db  db
	log log
}

func NewHandler(db db, log log) *Handler {
	return &Handler{
		db:  db,
		log: log,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, _ *http.Request) {
	if h.db == nil {
		h.log.Error("Instance DB isn't initialized")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := h.db.Ping()
	if err != nil {
		h.log.WithError(err).Error("Ping database is failed")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
