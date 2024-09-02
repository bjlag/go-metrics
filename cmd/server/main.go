package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/bjlag/go-metrics/internal/handler/update_common"
	"github.com/bjlag/go-metrics/internal/handler/update_counter"
	"github.com/bjlag/go-metrics/internal/handler/update_gauge"
	"github.com/bjlag/go-metrics/internal/helper"
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
	mux := http.NewServeMux()
	mux.HandleFunc("/update/", update_common.Handle)
	mux.Handle("/update/gauge/", helper.MakeUpdateHandler(update_gauge.Handle))
	mux.Handle("/update/counter/", helper.MakeUpdateHandler(update_counter.Handle))

	log.Printf("Listening on %s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		return err
	}

	return nil
}
