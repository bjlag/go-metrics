package main

import (
	nativLog "log"
	"net/http"

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

	memStorage := memory.NewStorage()
	htmlRenderer := renderer.NewHTMLRenderer(tmplPath)

	backup, err := file.NewStorage(fileStoragePath, 0)
	if err != nil {
		log.Error("Failed to create backup storage", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	data, err := backup.Load()
	if err != nil {
		log.Error("Failed to load backup data", map[string]interface{}{
			"error": err.Error(),
		})
	}

	for _, value := range data {
		switch value.MType {
		case "counter":
			memStorage.AddCounter(value.ID, *value.Delta)
		case "gauge":
			memStorage.SetGauge(value.ID, *value.Value)
		}
	}

	return http.ListenAndServe(addr.String(), initRouter(htmlRenderer, memStorage, backup, log))
}
