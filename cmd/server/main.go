package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bjlag/go-metrics/internal/handler/update"
	"github.com/bjlag/go-metrics/internal/handler/update_counter"
	"github.com/bjlag/go-metrics/internal/handler/update_gauge"
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
	gaugeHandler := update_gauge.NewHandler(memStorage)
	counterHandler := update_counter.NewHandler(memStorage)
	updateHandler := update.NewHandler()

	mux := http.NewServeMux()
	mux.Handle("/update/gauge/{name}/{value}", middleware.Default(http.HandlerFunc(gaugeHandler.Handle)))
	mux.Handle("/update/counter/{name}/{value}", middleware.Default(http.HandlerFunc(counterHandler.Handle)))
	mux.Handle("/update/", http.HandlerFunc(updateHandler.Handle))

	log.Printf("Listening on %s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		return err
	}

	return nil
}
