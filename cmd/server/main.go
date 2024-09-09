package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	updateCounter "github.com/bjlag/go-metrics/internal/handler/update/counter"
	updateGauge "github.com/bjlag/go-metrics/internal/handler/update/gauge"
	updateUnknown "github.com/bjlag/go-metrics/internal/handler/update/unknown"
	valueCaunter "github.com/bjlag/go-metrics/internal/handler/value/counter"
	valueGauge "github.com/bjlag/go-metrics/internal/handler/value/gauge"
	valueUnknown "github.com/bjlag/go-metrics/internal/handler/value/unknown"
	"github.com/bjlag/go-metrics/internal/middleware"
	storage "github.com/bjlag/go-metrics/internal/storage/memory"
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
	memStorage := storage.NewMemStorage()

	router := chi.NewRouter()

	router.Use(
		middleware.LogRequestMiddleware,
		middleware.FinishRequestMiddleware,
	)

	router.Post("/update/gauge/{name}/{value}", updateGauge.NewHandler(memStorage).Handle)
	router.Post("/update/counter/{name}/{value}", updateCounter.NewHandler(memStorage).Handle)
	router.Post("/update/{kind}/{name}/{value}", updateUnknown.NewHandler().Handle)

	router.Get("/value/gauge/{name}", valueGauge.NewHandler(memStorage).Handle)
	router.Get("/value/counter/{name}", valueCaunter.NewHandler(memStorage).Handle)
	router.Get("/value/{kind}/{name}", valueUnknown.NewHandler().Handle)

	log.Printf("Listening on %s", port)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
