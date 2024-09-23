package main

import (
	"github.com/go-chi/chi/v5"

	"github.com/bjlag/go-metrics/internal/handler/list"
	updateCounter "github.com/bjlag/go-metrics/internal/handler/update/counter"
	updateGauge "github.com/bjlag/go-metrics/internal/handler/update/gauge"
	updateGaneral "github.com/bjlag/go-metrics/internal/handler/update/general"
	updateUnknown "github.com/bjlag/go-metrics/internal/handler/update/unknown"
	valueCaunter "github.com/bjlag/go-metrics/internal/handler/value/counter"
	valueGauge "github.com/bjlag/go-metrics/internal/handler/value/gauge"
	valueGaneral "github.com/bjlag/go-metrics/internal/handler/value/general"
	valueUnknown "github.com/bjlag/go-metrics/internal/handler/value/unknown"
	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/middleware"
	"github.com/bjlag/go-metrics/internal/storage"
	"github.com/bjlag/go-metrics/internal/util/renderer"
)

func initRouter(htmlRenderer *renderer.HTMLRenderer, memStorage storage.Repository, logger logger.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		middleware.LogRequest(logger),
	)

	r.Get("/", list.NewHandler(htmlRenderer, memStorage).Handle)

	r.Route("/update", func(r chi.Router) {
		jsonContentType := middleware.SetHeaderResponse("Content-Type", []string{"application/json"})
		textContentType := middleware.SetHeaderResponse("Content-Type", []string{"text/plain", "charset=utf-8"})

		r.With(jsonContentType).Post("/", updateGaneral.NewHandler(memStorage).Handle)
		r.With(textContentType).Post("/gauge/{name}/{value}", updateGauge.NewHandler(memStorage).Handle)
		r.With(textContentType).Post("/counter/{name}/{value}", updateCounter.NewHandler(memStorage).Handle)
		r.With(textContentType).Post("/{kind}/{name}/{value}", updateUnknown.NewHandler().Handle)
	})

	r.Route("/value", func(r chi.Router) {
		jsonContentType := middleware.SetHeaderResponse("Content-Type", []string{"application/json"})
		textContentType := middleware.SetHeaderResponse("Content-Type", []string{"text/plain", "charset=utf-8"})

		r.With(jsonContentType).Post("/", valueGaneral.NewHandler(memStorage).Handle)
		r.With(textContentType).Post("/gauge/{name}", valueGauge.NewHandler(memStorage).Handle)
		r.With(textContentType).Post("/counter/{name}", valueCaunter.NewHandler(memStorage).Handle)
		r.With(textContentType).Post("/{kind}/{name}", valueUnknown.NewHandler().Handle)
	})

	return r
}
