package list

import (
	"net/http"

	"github.com/bjlag/go-metrics/internal/storage"
)

const (
	writeBodyMsgErr = "Error while writing body"
)

type Handler struct {
	renderer renderer
	repo     repo
	log      log
}

func NewHandler(renderer renderer, repo repo, log log) *Handler {
	return &Handler{
		renderer: renderer,
		repo:     repo,
		log:      log,
	}
}

func (h Handler) Handle(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Title    string
		Gauges   storage.Gauges
		Counters storage.Counters
	}{
		Title:    "Список метрик",
		Gauges:   h.repo.GetAllGauges(),
		Counters: h.repo.GetAllCounters(),
	}

	err := h.renderer.Render(w, "list.html", data)
	if err != nil {
		h.log.WithField("err", err.Error()).
			Error("failed to render list.html")
		http.Error(w, writeBodyMsgErr, http.StatusInternalServerError)
	}
}
