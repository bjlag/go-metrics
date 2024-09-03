package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bjlag/go-metrics/internal/handler/counter"
	"github.com/bjlag/go-metrics/internal/handler/gauge"
	"github.com/bjlag/go-metrics/internal/handler/update"
	"github.com/bjlag/go-metrics/internal/middleware"
	"github.com/bjlag/go-metrics/internal/storage/memory"
)

const (
	port = "8080"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	memStorage := memory.NewMemStorage()
	gaugeHandler := gauge.NewHandler(memStorage)
	counterHandler := counter.NewHandler(memStorage)
	updateHandler := update.NewHandler()

	mux := http.NewServeMux()
	mux.Handle("/update/gauge/{name}/{value}", middleware.Default(http.HandlerFunc(gaugeHandler.Handle)))
	mux.Handle("/update/counter/{name}/{value}", middleware.Default(http.HandlerFunc(counterHandler.Handle)))
	mux.Handle("/update/{type}/", http.HandlerFunc(updateHandler.Handle))

	log.Printf("Listening on %s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		return err
	}

	return nil
}
