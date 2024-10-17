package main

import (
	"github.com/bjlag/go-metrics/internal/signature"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"

	"github.com/bjlag/go-metrics/internal/backup"
	"github.com/bjlag/go-metrics/internal/handler/list"
	"github.com/bjlag/go-metrics/internal/handler/ping"
	updateBatch "github.com/bjlag/go-metrics/internal/handler/update/batch"
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
	"github.com/bjlag/go-metrics/internal/renderer"
	"github.com/bjlag/go-metrics/internal/storage"
)

func initRouter(htmlRenderer *renderer.HTMLRenderer, storage storage.Repository, db *sqlx.DB, backup backup.Creator, singManager *signature.SignManager, log logger.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		middleware.NewLogRequest(log).Handle,
		middleware.NewGzip(log).Handle,
	)

	r.Route("/", func(r chi.Router) {
		r.With(middleware.SetHeaderResponse("Content-Type", "text/html")).
			Get("/", list.NewHandler(htmlRenderer, storage, log).Handle)
	})

	r.Route("/update", func(r chi.Router) {
		jsonContentType := middleware.SetHeaderResponse("Content-Type", "application/json")
		textContentType := middleware.SetHeaderResponse("Content-Type", "text/plain", "charset=utf-8")

		r.With(jsonContentType).Post("/", updateGaneral.NewHandler(storage, backup, log).Handle)
		r.With(textContentType).Post("/gauge/{name}/{value}", updateGauge.NewHandler(storage, backup, log).Handle)
		r.With(textContentType).Post("/counter/{name}/{value}", updateCounter.NewHandler(storage, backup, log).Handle)
		r.With(textContentType).Post("/{kind}/{name}/{value}", updateUnknown.NewHandler(log).Handle)
	})

	r.Route("/updates", func(r chi.Router) {
		jsonContentType := middleware.SetHeaderResponse("Content-Type", "application/json")

		r.With(jsonContentType).Post("/", updateBatch.NewHandler(storage, singManager, backup, log).Handle)
	})

	r.Route("/value", func(r chi.Router) {
		jsonContentType := middleware.SetHeaderResponse("Content-Type", "application/json")
		textContentType := middleware.SetHeaderResponse("Content-Type", "text/plain", "charset=utf-8")

		r.With(jsonContentType).Post("/", valueGaneral.NewHandler(storage, log).Handle)
		r.With(textContentType).Get("/gauge/{name}", valueGauge.NewHandler(storage, log).Handle)
		r.With(textContentType).Get("/counter/{name}", valueCaunter.NewHandler(storage, log).Handle)
		r.With(textContentType).Get("/{kind}/{name}", valueUnknown.NewHandler(log).Handle)
	})

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", ping.NewHandler(db, log).Handle)
	})

	return r
}
