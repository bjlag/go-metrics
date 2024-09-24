package list

import (
	"fmt"
	"net/http"

	"github.com/bjlag/go-metrics/internal/storage"
)

const (
	writeBodyMsgErr = "Error while writing body"
)

type Handler struct {
	renderer Renderer
	storage  Storage
	log      Logger
}

func NewHandler(renderer Renderer, storage Storage, logger Logger) *Handler {
	return &Handler{
		renderer: renderer,
		storage:  storage,
		log:      logger,
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
		h.log.Error(fmt.Sprintf("Failed to render list.html: %s", err.Error()), nil)
		http.Error(w, writeBodyMsgErr, http.StatusInternalServerError)
	}
}
