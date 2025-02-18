package list

import (
	"net/http"

	"github.com/bjlag/go-metrics/internal/storage"
)

const (
	writeBodyMsgErr = "Error while writing body"
)

// Handler обработчик HTTP запроса для получения и вывода списка метрик на HTML странице.
type Handler struct {
	renderer renderer
	repo     repo
	log      log
}

// NewHandler создает обработчик.
func NewHandler(renderer renderer, repo repo, log log) *Handler {
	return &Handler{
		renderer: renderer,
		repo:     repo,
		log:      log,
	}
}

// Handle обрабатывает HTTP запрос.
func (h Handler) Handle(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title    string
		Gauges   storage.Gauges
		Counters storage.Counters
	}{
		Title:    "Список метрик",
		Gauges:   h.repo.GetAllGauges(r.Context()),
		Counters: h.repo.GetAllCounters(r.Context()),
	}

	err := h.renderer.Render(w, "list.html", data)
	if err != nil {
		h.log.WithError(err).Error("Failed to render list.html")
		http.Error(w, writeBodyMsgErr, http.StatusInternalServerError)
	}
}
