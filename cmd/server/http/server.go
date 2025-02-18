package http

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/sync/errgroup"

	"github.com/bjlag/go-metrics/internal/backup"
	asyncBackup "github.com/bjlag/go-metrics/internal/backup/async"
	syncBackup "github.com/bjlag/go-metrics/internal/backup/sync"
	"github.com/bjlag/go-metrics/internal/http/handler/list"
	"github.com/bjlag/go-metrics/internal/http/handler/ping"
	updateBatch "github.com/bjlag/go-metrics/internal/http/handler/update/batch"
	updateCounter "github.com/bjlag/go-metrics/internal/http/handler/update/counter"
	updateGauge "github.com/bjlag/go-metrics/internal/http/handler/update/gauge"
	updateGaneral "github.com/bjlag/go-metrics/internal/http/handler/update/general"
	updateUnknown "github.com/bjlag/go-metrics/internal/http/handler/update/unknown"
	valueCaunter "github.com/bjlag/go-metrics/internal/http/handler/value/counter"
	valueGauge "github.com/bjlag/go-metrics/internal/http/handler/value/gauge"
	valueGaneral "github.com/bjlag/go-metrics/internal/http/handler/value/general"
	valueUnknown "github.com/bjlag/go-metrics/internal/http/handler/value/unknown"
	middleware2 "github.com/bjlag/go-metrics/internal/http/middleware"
	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/renderer"
	"github.com/bjlag/go-metrics/internal/securety/crypt"
	"github.com/bjlag/go-metrics/internal/securety/signature"
	"github.com/bjlag/go-metrics/internal/storage"
)

const (
	gdTimeout = 10 * time.Second
)

type Server struct {
	addr          string
	htmlRenderer  *renderer.HTMLRenderer
	repo          storage.Repository
	db            *sqlx.DB
	backup        backup.Creator
	singManager   *signature.SignManager
	cryptManager  *crypt.DecryptManager
	trustedSubnet *net.IPNet
	log           logger.Logger
}

func NewServer(
	addr string,
	htmlRenderer *renderer.HTMLRenderer,
	repo storage.Repository,
	db *sqlx.DB,
	backup backup.Creator,
	singManager *signature.SignManager,
	cryptManager *crypt.DecryptManager,
	trustedSubnet *net.IPNet,
	log logger.Logger,
) *Server {
	return &Server{
		addr:          addr,
		htmlRenderer:  htmlRenderer,
		repo:          repo,
		db:            db,
		backup:        backup,
		singManager:   singManager,
		cryptManager:  cryptManager,
		trustedSubnet: trustedSubnet,
		log:           log,
	}
}

func (s *Server) Start(ctx context.Context) error {
	httpServer := &http.Server{
		Addr:    s.addr,
		Handler: s.router(),
	}

	s.log.Info("Starting HTTP server")

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()

		s.log.Info("Shutting down HTTP server")

		gdCtx, cancel := context.WithTimeout(context.Background(), gdTimeout)
		defer cancel()

		switch b := s.backup.(type) {
		case *asyncBackup.Backup:
			b.Stop(gdCtx)
		case *syncBackup.Backup:
			err := b.Create(gdCtx)
			if err != nil {
				s.log.WithError(err).Error("Failed to create backup while shutting down")
				return err
			}
		default:
			return errors.New("unknown backup type")
		}

		return httpServer.Shutdown(gdCtx)
	})

	if err := g.Wait(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	return nil
}

func (s *Server) router() http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware2.LogMiddleware(s.log),
		middleware2.CheckRealIPMiddleware(s.trustedSubnet, s.log),
		middleware2.GzipMiddleware(s.log),
		middleware2.DecryptMiddleware(s.cryptManager, s.log),
	)

	r.Route("/", func(r chi.Router) {
		r.With(middleware2.HeaderResponseMiddleware("Content-Type", "text/html")).
			Get("/", list.NewHandler(s.htmlRenderer, s.repo, s.log).Handle)
	})

	r.Route("/update", func(r chi.Router) {
		jsonContentType := middleware2.HeaderResponseMiddleware("Content-Type", "application/json")
		textContentType := middleware2.HeaderResponseMiddleware("Content-Type", "text/plain", "charset=utf-8")

		r.With(jsonContentType).Post("/", updateGaneral.NewHandler(s.repo, s.backup, s.log).Handle)
		r.With(textContentType).Post("/gauge/{name}/{value}", updateGauge.NewHandler(s.repo, s.backup, s.log).Handle)
		r.With(textContentType).Post("/counter/{name}/{value}", updateCounter.NewHandler(s.repo, s.backup, s.log).Handle)
		r.With(textContentType).Post("/{kind}/{name}/{value}", updateUnknown.NewHandler(s.log).Handle)
	})

	r.Route("/updates", func(r chi.Router) {
		jsonContentType := middleware2.HeaderResponseMiddleware("Content-Type", "application/json")
		validateSignRequest := middleware2.SignatureMiddleware(s.singManager, s.log)

		r.
			With(jsonContentType).
			With(validateSignRequest).
			Post("/", updateBatch.NewHandler(s.repo, s.backup, s.log).Handle)
	})

	r.Route("/value", func(r chi.Router) {
		jsonContentType := middleware2.HeaderResponseMiddleware("Content-Type", "application/json")
		textContentType := middleware2.HeaderResponseMiddleware("Content-Type", "text/plain", "charset=utf-8")

		r.With(jsonContentType).Post("/", valueGaneral.NewHandler(s.repo, s.log).Handle)
		r.With(textContentType).Get("/gauge/{name}", valueGauge.NewHandler(s.repo, s.log).Handle)
		r.With(textContentType).Get("/counter/{name}", valueCaunter.NewHandler(s.repo, s.log).Handle)
		r.With(textContentType).Get("/{kind}/{name}", valueUnknown.NewHandler(s.log).Handle)
	})

	r.Route("/ping", func(r chi.Router) {
		r.Get("/", ping.NewHandler(s.db, s.log).Handle)
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
