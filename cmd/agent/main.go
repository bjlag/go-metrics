package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second

	timeout = 100 * time.Millisecond

	urlFormat = "http://127.0.0.1:8080/update/%s/%s/%v"

	gaugeMetric   = "gauge"
	counterMetric = "counter"
)

func main() {
	log.Println("Starting agent")

	rtm := &runtime.MemStats{}

	client := &http.Client{}
	client.Timeout = timeout

	pollIntervalTicker := time.NewTicker(pollInterval)
	defer pollIntervalTicker.Stop()

	go func() {
		for ; ; <-pollIntervalTicker.C {
			runtime.ReadMemStats(rtm)

			response, err := sendMetric(client, NewMetric(counterMetric, "PollCount", 1))
			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("Sent request to %s, status %d", response.Request.URL.Path, response.StatusCode)
		}
	}()

	reportIntervalTicker := time.NewTicker(reportInterval)
	defer reportIntervalTicker.Stop()

	wg := &sync.WaitGroup{}

	for ; ; <-reportIntervalTicker.C {
		for _, metric := range collectMetrics(rtm) {
			wg.Add(1)
			go func() {
				defer wg.Done()

				response, err := sendMetric(client, metric)
				if err != nil {
					log.Println(err)
					return
				}

				defer func() {
					_ = response.Body.Close()
				}()

				log.Printf("Sent request to %s, status %d", response.Request.URL.Path, response.StatusCode)
			}()
		}

		wg.Wait()
	}
}

type Metric struct {
	kind  string
	name  string
	value interface{}
}

func NewMetric(kind string, name string, value interface{}) *Metric {
	return &Metric{
		kind:  kind,
		name:  name,
		value: value,
	}
}

func sendMetric(client *http.Client, metric *Metric) (*http.Response, error) {
	url := fmt.Sprintf(urlFormat, metric.kind, metric.name, metric.value)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request to '%s', error %v", url, err)
	}

	request.Header.Set("Content-Type", "text/plain")

	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request to '%s', error %v", url, err)
	}

	return response, nil
}

func collectMetrics(rtm *runtime.MemStats) []*Metric {
	return []*Metric{
		NewMetric(gaugeMetric, "Alloc", rtm.Alloc),
		NewMetric(gaugeMetric, "TotalAlloc", rtm.TotalAlloc),
		NewMetric(gaugeMetric, "BuckHashSys", rtm.BuckHashSys),
		NewMetric(gaugeMetric, "Frees", rtm.Frees),
		NewMetric(gaugeMetric, "GCCPUFraction", rtm.GCCPUFraction),
		NewMetric(gaugeMetric, "GCSys", rtm.GCSys),
		NewMetric(gaugeMetric, "HeapAlloc", rtm.HeapAlloc),
		NewMetric(gaugeMetric, "HeapIdle", rtm.HeapIdle),
		NewMetric(gaugeMetric, "HeapInuse", rtm.HeapInuse),
		NewMetric(gaugeMetric, "HeapObjects", rtm.HeapObjects),
		NewMetric(gaugeMetric, "HeapReleased", rtm.HeapReleased),
		NewMetric(gaugeMetric, "HeapSys", rtm.HeapSys),
		NewMetric(gaugeMetric, "LastGC", rtm.LastGC),
		NewMetric(gaugeMetric, "Lookups", rtm.Lookups),
		NewMetric(gaugeMetric, "MCacheInuse", rtm.MCacheInuse),
		NewMetric(gaugeMetric, "MCacheSys", rtm.MCacheSys),
		NewMetric(gaugeMetric, "MSpanInuse", rtm.MSpanInuse),
		NewMetric(gaugeMetric, "MSpanSys", rtm.MSpanSys),
		NewMetric(gaugeMetric, "Mallocs", rtm.Mallocs),
		NewMetric(gaugeMetric, "NextGC", rtm.NextGC),
		NewMetric(gaugeMetric, "NumForcedGC", rtm.NumForcedGC),
		NewMetric(gaugeMetric, "NumGC", rtm.NumGC),
		NewMetric(gaugeMetric, "OtherSys", rtm.OtherSys),
		NewMetric(gaugeMetric, "PauseTotalNs", rtm.PauseTotalNs),
		NewMetric(gaugeMetric, "StackInuse", rtm.StackInuse),
		NewMetric(gaugeMetric, "StackSys", rtm.StackSys),
		NewMetric(gaugeMetric, "Sys", rtm.Sys),
		NewMetric(gaugeMetric, "RandomValue", getRandomFloat(1, 100)),
	}
}

func getRandomFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
