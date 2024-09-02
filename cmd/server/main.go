package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
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
	mux.HandleFunc("/update/", updateCommonHandler)
	mux.Handle("/update/gauge/", makeUpdateHandler(updateGaugeHandler))
	mux.Handle("/update/counter/", makeUpdateHandler(updateCounterHandler))

	log.Printf("Listening on %s", port)

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
	if err != nil {
		return err
	}

	return nil
}

const (
	invalidMetricTypeMsgErr = "Invalid metric type"
	noNameMetricMsgErr      = "Metric name not specified"
	invalidMetricPathMsgErr = "Invalid metric path"
	invalidTypeValueMsgErr  = "Invalid type value of metric"
)

var (
	validRoutePattern  = regexp.MustCompile("^/update/(gauge|counter)/([a-zA-Z0-9_]+)?/(\\d+(.\\d+)?)$")
	withoutNamePattern = regexp.MustCompile("^/update/(gauge|counter)/(\\d+(.\\d+)?)$")
)

func updateCommonHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, invalidMetricTypeMsgErr, http.StatusBadRequest)
}

func updateGaugeHandler(w http.ResponseWriter, r *http.Request, nameMetric, valueMetric string) {
	value, err := strconv.ParseFloat(valueMetric, 64)
	if err != nil {
		http.Error(w, invalidTypeValueMsgErr, http.StatusBadRequest)
		return
	}

	_ = value
}

func updateCounterHandler(w http.ResponseWriter, r *http.Request, nameMetric, valueMetric string) {
	value, err := strconv.ParseInt(valueMetric, 10, 64)
	if err != nil {
		http.Error(w, invalidTypeValueMsgErr, http.StatusBadRequest)
		return
	}

	_ = value
}

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
