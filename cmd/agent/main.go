package main

import (
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/agent/sender"
)

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second

	timeout = 100 * time.Millisecond
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	log.Println("Starting agent")

	metricCollector := collector.NewMetricCollector(&runtime.MemStats{})

	client := &http.Client{}
	client.Timeout = timeout
	metricSender := sender.NewHttpSender(client)

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()

	go func() {
		for ; ; <-pollTicker.C {
			metricCollector.ReadStats()

			response, err := metricSender.Send(collector.NewMetric(collector.Counter, "PollCount", 1))
			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("Sent request to %s, status %d", response.Request.URL.Path, response.StatusCode)
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

				response, err := metricSender.Send(metric)
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
