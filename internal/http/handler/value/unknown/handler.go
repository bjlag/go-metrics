package unknown

import (
	"net/http"
)

// Handler обработчик HTTP запроса на случай если получаем значение метрики неизвестного типа.
type Handler struct {
	log log
}

// NewHandler создает обработчик.
func NewHandler(log log) *Handler {
	return &Handler{
		log: log,
	}
}

// Handle обрабатывает HTTP запрос.
func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	h.log.WithField("type", r.PathValue("kind")).
		WithField("url", r.URL.Path).
		Info("Invalid metric type")
	http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
}
