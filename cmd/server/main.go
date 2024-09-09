package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

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

	router := chi.NewRouter()

	router.Use(
		middleware.LogRequestMiddleware,
		middleware.AllowPostMethodMiddleware,
		middleware.FinishRequestMiddleware,
	)

	router.Post("/update/gauge/{name}/{value}", gauge.NewHandler(memStorage).Handle)
	router.Post("/update/counter/{name}/{value}", counter.NewHandler(memStorage).Handle)
	router.Post("/update/{kind}/{name}/{value}", update.NewHandler().Handle)

	log.Printf("Listening on %s", port)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), router)
}
