package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/bjlag/go-metrics/internal/agent/client"
	"github.com/bjlag/go-metrics/internal/agent/collector"
)

const (
	baseURL = "http://127.0.0.1:8080"

	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	log.Println("Starting agent")

	metricCollector := collector.NewMetricCollector(&runtime.MemStats{})
	metricClient := client.NewHTTPSender(baseURL)

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()

	go func() {
		for ; ; <-pollTicker.C {
			metricCollector.ReadStats()

			response, err := metricClient.Send(collector.NewMetric(collector.Counter, "PollCount", 1))
			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("Sent request to %s, status %d", response.Request.URL, response.StatusCode())
		}
	}()

	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	wg := &sync.WaitGroup{}

	for ; ; <-reportTicker.C {
		for _, metric := range metricCollector.Collect() {
			wg.Add(1)
			go func() {
				defer wg.Done()

				response, err := metricClient.Send(metric)
				if err != nil {
					log.Println(err)
					return
				}

				log.Printf("Sent request to %s, status %d", response.Request.URL, response.StatusCode())
			}()
		}

		wg.Wait()
	}
}
