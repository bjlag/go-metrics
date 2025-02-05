package main

import (
	"net"
	"net/http"
	"net/http/pprof"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"

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
	"github.com/bjlag/go-metrics/internal/securety/crypt"
	"github.com/bjlag/go-metrics/internal/securety/signature"
	"github.com/bjlag/go-metrics/internal/storage"
)

func initRouter(
	htmlRenderer *renderer.HTMLRenderer,
	storage storage.Repository,
	db *sqlx.DB,
	backup backup.Creator,
	singManager *signature.SignManager,
	crypt *crypt.DecryptManager,
	trustedSubnet *net.IPNet,
	log logger.Logger,
) *chi.Mux {
	r := chi.NewRouter()

	r.Use(
		middleware.LogMiddleware(log),
		middleware.CheckRealIPMiddleware(trustedSubnet, log),
		middleware.GzipMiddleware(log),
		middleware.DecryptMiddleware(crypt, log),
	)

	r.Route("/", func(r chi.Router) {
		r.With(middleware.HeaderResponseMiddleware("Content-Type", "text/html")).
			Get("/", list.NewHandler(htmlRenderer, storage, log).Handle)
	})

	r.Route("/update", func(r chi.Router) {
		jsonContentType := middleware.HeaderResponseMiddleware("Content-Type", "application/json")
		textContentType := middleware.HeaderResponseMiddleware("Content-Type", "text/plain", "charset=utf-8")

		r.With(jsonContentType).Post("/", updateGaneral.NewHandler(storage, backup, log).Handle)
		r.With(textContentType).Post("/gauge/{name}/{value}", updateGauge.NewHandler(storage, backup, log).Handle)
		r.With(textContentType).Post("/counter/{name}/{value}", updateCounter.NewHandler(storage, backup, log).Handle)
		r.With(textContentType).Post("/{kind}/{name}/{value}", updateUnknown.NewHandler(log).Handle)
	})

	r.Route("/updates", func(r chi.Router) {
		jsonContentType := middleware.HeaderResponseMiddleware("Content-Type", "application/json")
		validateSignRequest := middleware.SignatureMiddleware(singManager, log)

		r.
			With(jsonContentType).
			With(validateSignRequest).
			Post("/", updateBatch.NewHandler(storage, backup, log).Handle)
	})

	r.Route("/value", func(r chi.Router) {
		jsonContentType := middleware.HeaderResponseMiddleware("Content-Type", "application/json")
		textContentType := middleware.HeaderResponseMiddleware("Content-Type", "text/plain", "charset=utf-8")

		r.With(jsonContentType).Post("/", valueGaneral.NewHandler(storage, log).Handle)
		r.With(textContentType).Get("/gauge/{name}", valueGauge.NewHandler(storage, log).Handle)
		r.With(textContentType).Get("/counter/{name}", valueCaunter.NewHandler(storage, log).Handle)
		r.With(textContentType).Get("/{kind}/{name}", valueUnknown.NewHandler(log).Handle)
	})

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", ping.NewHandler(db, log).Handle)
	})

	r.Route("/debug/pprof", func(r chi.Router) {
		r.Use(chiMiddleware.NoCache)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			if r.RequestURI[len(r.RequestURI)-1:] != "/" {
				http.Redirect(w, r, r.RequestURI+"/", http.StatusMovedPermanently)
			}

			pprof.Index(w, r)
		})

		r.Get("/cmdline", pprof.Cmdline)
		r.Get("/profile", pprof.Profile)
		r.Get("/symbol", pprof.Symbol)
		r.Get("/trace", pprof.Trace)

		r.Get("/goroutine", pprof.Handler("goroutine").ServeHTTP)
		r.Get("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
		r.Get("/mutex", pprof.Handler("mutex").ServeHTTP)
		r.Get("/heap", pprof.Handler("heap").ServeHTTP)
		r.Get("/block", pprof.Handler("block").ServeHTTP)
		r.Get("/allocs", pprof.Handler("allocs").ServeHTTP)
	})

	r.Route("/docs", func(r chi.Router) {
		r.Get("/*", httpSwagger.Handler())
	})

	return r
}
