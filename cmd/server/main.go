package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	updateCounter "github.com/bjlag/go-metrics/internal/handler/update/counter"
	updateGauge "github.com/bjlag/go-metrics/internal/handler/update/gauge"
	updateUnknown "github.com/bjlag/go-metrics/internal/handler/update/unknown"
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

	router := chi.NewRouter()

	router.Use(
		middleware.LogRequestMiddleware,
		middleware.FinishRequestMiddleware,
	)

	router.Post("/update/gauge/{name}/{value}", updateGauge.NewHandler(memStorage).Handle)
	router.Post("/update/counter/{name}/{value}", updateCounter.NewHandler(memStorage).Handle)
	router.Post("/update/{kind}/{name}/{value}", updateUnknown.NewHandler().Handle)

	log.Printf("Listening on %s", port)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
