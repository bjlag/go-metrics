package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/bjlag/go-metrics/internal/handler/update_common"
	"github.com/bjlag/go-metrics/internal/handler/update_counter"
	"github.com/bjlag/go-metrics/internal/handler/update_gauge"
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
	mux.Handle("/update/gauge/", makeUpdateHandler(update_gauge.Handle))
	mux.Handle("/update/counter/", makeUpdateHandler(update_counter.Handle))

	log.Printf("Listening on %s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		return err
	}

	return nil
}

const (
	noNameMetricMsgErr      = "Metric name not specified"
	invalidMetricPathMsgErr = "Invalid metric path"
)

var (
	validRoutePattern  = regexp.MustCompile("^/update/(gauge|counter)/([a-zA-Z0-9_]+)?/(\\d+(.\\d+)?)$")
	withoutNamePattern = regexp.MustCompile("^/update/(gauge|counter)/(\\d+(.\\d+)?)$")
)

func makeUpdateHandler(handler func(w http.ResponseWriter, r *http.Request, name, value string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}

		if withoutNamePattern.MatchString(r.URL.Path) {
			http.Error(w, noNameMetricMsgErr, http.StatusNotFound)
			return
		}

		m := validRoutePattern.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.Error(w, invalidMetricPathMsgErr, http.StatusBadRequest)
			return
		}

		log.Printf("Metric received: type '%s', name '%s', value '%s'\n", m[1], m[2], m[3])

		handler(w, r, m[2], m[3])
	}
}
