package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/bjlag/go-metrics/internal/agent/client"
	"github.com/bjlag/go-metrics/internal/agent/collector"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	parseFlags()

	log.Println("Starting agent")
	log.Printf("Sending metrics to %s\n", addr.String())
	log.Printf("Poll interval %s\n", pollInterval)
	log.Printf("Report interval %s\n", reportInterval)

	metricCollector := collector.NewMetricCollector(&runtime.MemStats{})
	metricClient := client.NewHTTPSender(addr.host, addr.port)

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
