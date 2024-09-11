package main

import (
	"log"
	"net/http"

	storage "github.com/bjlag/go-metrics/internal/storage/memory"
	"github.com/bjlag/go-metrics/internal/util/renderer"
)

const (
	tmplPath = "web/tmpl/list.html"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	parseFlags()
	parseEnvs()

	log.Printf("Listening on %s\n", addr.String())

	memStorage := storage.NewMemStorage()
	htmlRenderer := renderer.NewHTMLRenderer(tmplPath)

	return http.ListenAndServe(addr.String(), initRouter(htmlRenderer, memStorage))
}
