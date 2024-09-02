package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bjlag/go-metrics/internal/handler/update_common"
	"github.com/bjlag/go-metrics/internal/handler/update_counter"
	"github.com/bjlag/go-metrics/internal/handler/update_gauge"
	"github.com/bjlag/go-metrics/internal/helper"
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
	gauge := update_gauge.NewHandler(memStorage)
	counter := update_counter.NewHandler(memStorage)

	mux := http.NewServeMux()
	mux.Handle("/update/gauge/", helper.MakeUpdateHandler(gauge.Handle))
	mux.Handle("/update/counter/", helper.MakeUpdateHandler(counter.Handle))
	mux.HandleFunc("/update/", update_common.Handle)

	log.Printf("Listening on %s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		return err
	}

	return nil
}
