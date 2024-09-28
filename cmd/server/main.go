package main

import (
	"fmt"
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
	log.Info(fmt.Sprintf("Log level '%s'", logLevel), nil)
	log.Info(fmt.Sprintf("Store interval %ds", storeInterval), nil)
	log.Info(fmt.Sprintf("File storage path '%s'", fileStoragePath), nil)
	log.Info(fmt.Sprintf("Restore metrics %v", restore), nil)

	memStorage := memory.NewStorage()
	htmlRenderer := renderer.NewHTMLRenderer(tmplPath)

	backup, err := file.NewStorage(fileStoragePath, storeInterval)
	if err != nil {
		log.Error("Failed to create backup storage", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}

	if restore {
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

		log.Info("Backup loaded", map[string]interface{}{})
	}

	return http.ListenAndServe(addr.String(), initRouter(htmlRenderer, memStorage, backup, log))
}
