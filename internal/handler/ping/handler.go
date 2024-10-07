package ping

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

type Handler struct {
	dsn string
	log log
}

func NewHandler(dsn string, log log) *Handler {
	return &Handler{
		dsn: dsn,
		log: log,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, _ *http.Request) {
	db, err := sqlx.Connect("pgx", h.dsn)
	if err != nil {
		h.log.WithError(err).Error("Unable to connect to database")
		w.WriteHeader(http.StatusInternalServerError)
	}

	err = db.Ping()
	if err != nil {
		h.log.WithError(err).Error("Ping database is failed")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
