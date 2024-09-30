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
	"github.com/bjlag/go-metrics/internal/storage/file"
	"github.com/bjlag/go-metrics/internal/storage/memory"
	"github.com/bjlag/go-metrics/internal/util/renderer"
)

const (
	tmplPath = "web/tmpl/list.html"
)

func main() {
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

	parseFlags()
	parseEnvs()

	log, err := logger.NewZapLogger(logLevel)
	if err != nil {
		return err
	}
	defer func() {
		_ = log.Close()
	}()

	log.Info("Started server", map[string]interface{}{
		"address": addr.String(),
	})
	log.Info(fmt.Sprintf("Log level '%s'", logLevel), nil)
	log.Info(fmt.Sprintf("Store interval %s", storeInterval), nil)
	log.Info(fmt.Sprintf("File storage path '%s'", fileStoragePath), nil)
	log.Info(fmt.Sprintf("Restore metrics %v", restore), nil)

	memStorage := memory.NewStorage()
	htmlRenderer := renderer.NewHTMLRenderer(tmplPath)

	fileStorage, err := file.NewStorage(fileStoragePath)
	if err != nil {
		log.Error("Failed to create file storage", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	if restore {
		err := restoreData(fileStorage, memStorage)
		if err != nil {
			log.Error("Failed to load backup data", map[string]interface{}{
				"error": err.Error(),
			})
		}

		log.Info("Backup loaded", nil)
	}

	var (
		b  backup.Creator
		ba *asyncBackup.Backup
	)

	if storeInterval <= 0 {
		b = syncBackup.New(memStorage, fileStorage, log)
	} else {
		ba = asyncBackup.New(memStorage, fileStorage, storeInterval, log)
		ba.Start()
		b = ba
	}

	httpServer := &http.Server{
		Addr:    addr.String(),
		Handler: initRouter(htmlRenderer, memStorage, b, log),
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()

		log.Info("Graceful shutting down server", nil)
		ba.Stop()
		return httpServer.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
