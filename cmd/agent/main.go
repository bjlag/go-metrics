package main

import (
	"fmt"
	logNativ "log"
	"runtime"
	"sync"
	"time"

	"github.com/bjlag/go-metrics/internal/agent/client"
	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/logger"
)

func main() {
	if err := run(); err != nil {
		logNativ.Fatalln(err)
	}
}

func run() error {
	parseFlags()
	parseEnvs()

	log, err := logger.NewZapLogger(logLevel)
	if err != nil {
		return err
	}
	defer func() {
		_ = log.Close()
	}()

	log.Info("Starting agent", nil)
	log.Info(fmt.Sprintf("Sending metrics to %s", addr.String()), nil)
	log.Info(fmt.Sprintf("Poll interval %s", pollInterval), nil)
	log.Info(fmt.Sprintf("Report interval %s", reportInterval), nil)
	log.Info(fmt.Sprintf("Log level '%s'", logLevel), nil)

	metricCollector := collector.NewMetricCollector(&runtime.MemStats{})
	metricClient := client.NewHTTPSender(addr.host, addr.port)

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()

	go func() {
		for range pollTicker.C {
			metricCollector.ReadStats()

			response, err := metricClient.Send(collector.NewMetric(collector.Counter, "PollCount", 1))
			if err != nil {
				log.Error(err.Error(), nil)
				continue
			}

			log.Info("Sent request", map[string]interface{}{
				"uri":      response.Request.URL,
				"response": string(response.Body()),
				"status":   response.StatusCode(),
			})
		}
	}()

	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	wg := &sync.WaitGroup{}

	for range reportTicker.C {
		for _, metric := range metricCollector.Collect() {
			wg.Add(1)
			go func() {
				defer wg.Done()

				response, err := metricClient.Send(metric)
				if err != nil {
					log.Error(err.Error(), nil)
					return
				}

				log.Info("Sent request", map[string]interface{}{
					"uri":      response.Request.URL,
					"response": string(response.Body()),
					"status":   response.StatusCode(),
				})
			}()
		}

		wg.Wait()
	}

	return nil
}
