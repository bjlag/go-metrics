package main

import (
	nativLog "log"
	"net/http"

	"github.com/bjlag/go-metrics/internal/logger"
	storage "github.com/bjlag/go-metrics/internal/storage/memory"
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

	memStorage := storage.NewMemStorage()
	htmlRenderer := renderer.NewHTMLRenderer(tmplPath)

	return http.ListenAndServe(addr.String(), initRouter(htmlRenderer, memStorage, log))
}
