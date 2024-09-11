package main

import (
	"github.com/go-chi/chi/v5"

	"github.com/bjlag/go-metrics/internal/handler/list"
	updateCounter "github.com/bjlag/go-metrics/internal/handler/update/counter"
	updateGauge "github.com/bjlag/go-metrics/internal/handler/update/gauge"
	updateUnknown "github.com/bjlag/go-metrics/internal/handler/update/unknown"
	valueCaunter "github.com/bjlag/go-metrics/internal/handler/value/counter"
	valueGauge "github.com/bjlag/go-metrics/internal/handler/value/gauge"
	valueUnknown "github.com/bjlag/go-metrics/internal/handler/value/unknown"
	"github.com/bjlag/go-metrics/internal/middleware"
	"github.com/bjlag/go-metrics/internal/storage"
	"github.com/bjlag/go-metrics/internal/util/renderer"
)

func initRouter(htmlRenderer *renderer.HTMLRenderer, memStorage storage.Repository) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		middleware.LogRequestMiddleware,
		middleware.FinishRequestMiddleware,
	)

	router.Get("/", list.NewHandler(htmlRenderer, memStorage).Handle)

	router.Post("/update/gauge/{name}/{value}", updateGauge.NewHandler(memStorage).Handle)
	router.Post("/update/counter/{name}/{value}", updateCounter.NewHandler(memStorage).Handle)
	router.Post("/update/{kind}/{name}/{value}", updateUnknown.NewHandler().Handle)

	router.Get("/value/gauge/{name}", valueGauge.NewHandler(memStorage).Handle)
	router.Get("/value/counter/{name}", valueCaunter.NewHandler(memStorage).Handle)
	router.Get("/value/{kind}/{name}", valueUnknown.NewHandler().Handle)

	return router
}
