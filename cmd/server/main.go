package main

import (
	"context"
	"fmt"
	nativLog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/bjlag/go-metrics/cmd"
	_ "github.com/bjlag/go-metrics/docs"
	"github.com/bjlag/go-metrics/internal/backup"
	asyncBackup "github.com/bjlag/go-metrics/internal/backup/async"
	syncBackup "github.com/bjlag/go-metrics/internal/backup/sync"
	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/renderer"
	"github.com/bjlag/go-metrics/internal/signature"
	"github.com/bjlag/go-metrics/internal/storage"
	"github.com/bjlag/go-metrics/internal/storage/file"
	"github.com/bjlag/go-metrics/internal/storage/memory"
	"github.com/bjlag/go-metrics/internal/storage/pg"
)

const (
	tmplPath = "web/tmpl/list.html"
	noValue  = "N/A"
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
	parseFlags()
	parseEnvs()

	log, err := logger.NewZapLog(logLevel)
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

	log.WithField("address", addr.String()).Info("Starting server")
	log.Info(fmt.Sprintf("Log level '%s'", logLevel))
	log.Info(fmt.Sprintf("Store interval %s", storeInterval))
	log.Info(fmt.Sprintf("File storage path '%s'", fileStoragePath))
	log.Info(fmt.Sprintf("Restore metrics %v", restore))

	if err := run(log); err != nil {
		log.WithError(err).Error("Error running server")
		nativLog.Fatalln(err)
	}
}

func run(log logger.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	db := initDB(databaseDSN, log)

	var store storage.Repository
	if db != nil {
		store = pg.NewStorage(db, log)
	} else {
		store = memory.NewStorage()
	}

	backupStore, err := file.NewStorage(fileStoragePath)
	if err != nil {
		log.WithError(err).Error("Failed to create file storage")
		return err
	}

	if restore {
		err := restoreData(ctx, backupStore, store)
		if err != nil {
			log.WithError(err).Error("Failed to load backup data")
		}

		log.Info("Backup loaded")
	}

	var (
		backupCreator      backup.Creator
		asyncBackupCreator *asyncBackup.Backup
	)

	if storeInterval <= 0 {
		backupCreator = syncBackup.New(store, backupStore, log)
	} else {
		asyncBackupCreator = asyncBackup.New(store, backupStore, storeInterval, log)
		asyncBackupCreator.Start(ctx)
		backupCreator = asyncBackupCreator
	}

	signManager := signature.NewSignManager(secretKey)
	htmlRenderer := renderer.NewHTMLRenderer(tmplPath)
	httpServer := &http.Server{
		Addr:    addr.String(),
		Handler: initRouter(htmlRenderer, store, db, backupCreator, signManager, log),
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()

		log.Info("Graceful shutting down server")
		asyncBackupCreator.Stop(ctx)
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
