package main

import (
	"context"
	"errors"
	"fmt"
	nativLog "log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/bjlag/go-metrics/cmd"
	"github.com/bjlag/go-metrics/cmd/server/config"
	_ "github.com/bjlag/go-metrics/docs"
	"github.com/bjlag/go-metrics/internal/backup"
	asyncBackup "github.com/bjlag/go-metrics/internal/backup/async"
	syncBackup "github.com/bjlag/go-metrics/internal/backup/sync"
	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/renderer"
	"github.com/bjlag/go-metrics/internal/securety/crypt"
	"github.com/bjlag/go-metrics/internal/securety/signature"
	"github.com/bjlag/go-metrics/internal/storage"
	"github.com/bjlag/go-metrics/internal/storage/file"
	"github.com/bjlag/go-metrics/internal/storage/memory"
	"github.com/bjlag/go-metrics/internal/storage/pg"
)

const (
	tmplPath  = "web/tmpl/list.html"
	gdTimeout = 10 * time.Second
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

//	@title			Go Metrics
//	@version		1.0
//	@description	Сервис сбора метрик и алертинга

func main() {
	cfg := config.LoadConfig()

	log, err := logger.NewZapLog(cfg.LogLevel)
	if err != nil {
		nativLog.Fatalln(err)
	}
	defer func() {
		_ = log.Close()
	}()

	build := cmd.NewBuild(buildVersion, buildDate, buildCommit)
	log.Info(build.VersionString())
	log.Info(build.DateString())
	log.Info(build.CommitString())

	log.WithField("address", cfg.Address.String()).Info("Starting server")
	log.Info(fmt.Sprintf("Log level '%s'", cfg.LogLevel))
	log.Info(fmt.Sprintf("Store interval %s", cfg.StoreInterval))
	log.Info(fmt.Sprintf("File storage path '%s'", cfg.FileStoragePath))
	log.Info(fmt.Sprintf("Restore metrics %v", cfg.Restore))
	log.Info(fmt.Sprintf("Private key %s", cfg.CryptoKeyPath))
	log.Info(fmt.Sprintf("JSON config %s", cfg.ConfigPath))

	if err := run(log, cfg); err != nil {
		log.WithError(err).Error("Error running server")
		nativLog.Fatalln(err)
	}
}

func run(log logger.Logger, cfg *config.Configuration) error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	db := initDB(cfg.DatabaseDSN, log)

	var store storage.Repository
	if db != nil {
		store = pg.NewStorage(db, log)
	} else {
		store = memory.NewStorage()
	}

	backupStore, err := file.NewStorage(cfg.FileStoragePath)
	if err != nil {
		log.WithError(err).Error("Failed to create file storage")
		return err
	}

	if cfg.Restore {
		err = restoreData(ctx, backupStore, store)
		if err != nil {
			log.WithError(err).Error("Failed to load backup data")
		}

		log.Info("Backup loaded")
	}

	var (
		backupCreator      backup.Creator
		asyncBackupCreator *asyncBackup.Backup
	)

	if cfg.StoreInterval <= 0 {
		backupCreator = syncBackup.New(store, backupStore, log)
	} else {
		asyncBackupCreator = asyncBackup.New(store, backupStore, cfg.StoreInterval, log)
		asyncBackupCreator.Start(ctx)
		backupCreator = asyncBackupCreator
	}

	cryptManager, err := crypt.NewDecryptManager(cfg.CryptoKeyPath)
	if err != nil {
		return err
	}

	signManager := signature.NewSignManager(cfg.SecretKey)
	htmlRenderer := renderer.NewHTMLRenderer(tmplPath)
	httpServer := &http.Server{
		Addr:    cfg.Address.String(),
		Handler: initRouter(htmlRenderer, store, db, backupCreator, signManager, cryptManager, cfg.TrustedSubnet, log),
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()

		gdCtx, cancel := context.WithTimeout(context.Background(), gdTimeout)
		defer cancel()

		if asyncBackupCreator != nil {
			asyncBackupCreator.Stop(gdCtx)
		} else {
			err = backupCreator.Create(gdCtx)
			if err != nil {
				log.WithError(err).Error("Failed to create backup while shutting down")
				return err
			}
		}

		return httpServer.Shutdown(gdCtx)
	})

	if err := g.Wait(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	}

	log.Info("Server shutdown gracefully")

	return nil
}
