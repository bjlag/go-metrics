package list

import (
	"net/http"

	"github.com/bjlag/go-metrics/internal/storage"
)

const (
	writeBodyMsgErr = "Error while writing body"
)

type Handler struct {
	renderer Renderer
	storage  Storage
}

func NewHandler(renderer Renderer, storage Storage) *Handler {
	return &Handler{
		renderer: renderer,
		storage:  storage,
	}
}

func (h Handler) Handle(w http.ResponseWriter, _ *http.Request) {
	data := struct {
		Title    string
		Gauges   storage.Gauges
		Counters storage.Counters
	}{
		Title:    "Список метрик",
		Gauges:   h.storage.GetAllGauges(),
		Counters: h.storage.GetAllCounters(),
	}

	err := h.renderer.Render(w, "list.html", data)
	if err != nil {
		http.Error(w, writeBodyMsgErr, http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
