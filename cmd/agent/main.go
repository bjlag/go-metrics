package main

import (
	"context"
	"fmt"
	logNativ "log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

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
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

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
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	g, gCtx := errgroup.WithContext(ctx)

	go func() {
		<-ctx.Done()

		log.Info("Graceful shutting down agent", nil)
	}()

	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				log.Info("Stopped read stats", nil)
				return nil
			case <-pollTicker.C:
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
		}
	})

	g.Go(func() error {
		gr, _ := errgroup.WithContext(ctx)

		for {
			select {
			case <-gCtx.Done():
				log.Info("Stopped send metrics", nil)
				return nil
			case <-reportTicker.C:
				for _, metric := range metricCollector.Collect() {
					gr.Go(func() error {
						select {
						case <-gCtx.Done():
							return nil
						default:
							response, err := metricClient.Send(metric)
							if err != nil {
								return err
							}

							log.Info("Sent request", map[string]interface{}{
								"uri":      response.Request.URL,
								"response": string(response.Body()),
								"status":   response.StatusCode(),
							})
						}

						return nil
					})
				}

				if err := gr.Wait(); err != nil {
					log.Error(err.Error(), nil)
				}
			}
		}
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
