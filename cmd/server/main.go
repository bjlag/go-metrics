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

	"github.com/bjlag/go-metrics/internal/backup"
	asyncBackup "github.com/bjlag/go-metrics/internal/backup/async"
	syncBackup "github.com/bjlag/go-metrics/internal/backup/sync"
	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/renderer"
	"github.com/bjlag/go-metrics/internal/storage"
	"github.com/bjlag/go-metrics/internal/storage/file"
	"github.com/bjlag/go-metrics/internal/storage/memory"
	"github.com/bjlag/go-metrics/internal/storage/pg"
)

const (
	tmplPath = "web/tmpl/list.html"
)

func main() {
	parseFlags()
	parseEnvs()

	if err := run(); err != nil {
		nativLog.Fatalln(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	log, err := logger.NewZapLog(logLevel)
	if err != nil {
		return err
	}
	defer func() {
		_ = log.Close()
	}()

	log.WithField("address", addr.String()).Info("started server")
	log.Info(fmt.Sprintf("log level '%s'", logLevel))
	log.Info(fmt.Sprintf("store interval %s", storeInterval))
	log.Info(fmt.Sprintf("file storage path '%s'", fileStoragePath))
	log.Info(fmt.Sprintf("restore metrics %v", restore))

	db := initDB(databaseDSN, log)

	var store storage.Repository
	if db != nil {
		store = pg.NewStorage(db, log)
	} else {
		store = memory.NewStorage()
	}

	backupStore, err := file.NewStorage(fileStoragePath)
	if err != nil {
		log.WithField("error", err.Error()).Error("failed to create file storage")
		return err
	}

	if restore {
		err := restoreData(ctx, backupStore, store)
		if err != nil {
			log.WithField("error", err.Error()).
				Error("failed to load backup data")
		}

		log.Info("backup loaded")
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

	htmlRenderer := renderer.NewHTMLRenderer(tmplPath)
	httpServer := &http.Server{
		Addr:    addr.String(),
		Handler: initRouter(htmlRenderer, store, db, backupCreator, log),
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()

		log.Info("graceful shutting down server")
		asyncBackupCreator.Stop(ctx)
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
