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
	// todo middleware на проверку метода на POST
	// todo middleware проверка заголовка Content-Type: text/plain
	// todo валидатор http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>
	// todo валидация значения метрики gauge float64, counter int64
	// todo редиректы не поддерживаются

	mux := http.NewServeMux()
	mux.HandleFunc("/update/", updateCommonHandler)
	mux.HandleFunc("/update/gauge/", updateGaugeHandler)
	mux.HandleFunc("/update/counter/", updateCounterHandler)

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

func updateGaugeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
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

	typeMetric := m[1]
	nameMetric := m[2]
	valueMetric, err := strconv.ParseFloat(m[3], 64)
	if err != nil {
		http.Error(w, invalidTypeValueMsgErr, http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(w, "type: %s\n", typeMetric)
	_, _ = fmt.Fprintf(w, "name: %s\n", nameMetric)
	_, _ = fmt.Fprintf(w, "value: %v\n", valueMetric)
}

func updateCounterHandler(w http.ResponseWriter, r *http.Request) {
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

	typeMetric := m[1]
	nameMetric := m[2]
	valueMetric, err := strconv.ParseInt(m[3], 10, 64)
	if err != nil {
		http.Error(w, invalidTypeValueMsgErr, http.StatusBadRequest)
		return
	}

	_, _ = fmt.Fprintf(w, "type: %s\n", typeMetric)
	_, _ = fmt.Fprintf(w, "name: %s\n", nameMetric)
	_, _ = fmt.Fprintf(w, "value: %v\n", valueMetric)
}
