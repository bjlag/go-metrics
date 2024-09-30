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

	log.Info("starting agent")
	log.Info(fmt.Sprintf("sending metrics to %s", addr.String()))
	log.Info(fmt.Sprintf("poll interval %s", pollInterval))
	log.Info(fmt.Sprintf("report interval %s", reportInterval))
	log.Info(fmt.Sprintf("log level '%s'", logLevel))

	metricCollector := collector.NewMetricCollector(&runtime.MemStats{})
	metricClient := client.NewHTTPSender(addr.host, addr.port)

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	g, gCtx := errgroup.WithContext(ctx)

	go func() {
		<-ctx.Done()

		log.Info("graceful shutting down agent")
	}()

	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				log.Info("stopped read stats")
				return nil
			case <-pollTicker.C:
				metricCollector.ReadStats()

				response, err := metricClient.Send(collector.NewMetric(collector.Counter, "PollCount", 1))
				if err != nil {
					log.WithField("error", err.Error()).
						Error("error in sending metric")
					continue
				}

				log.WithField("uri", response.Request.URL).
					WithField("response", string(response.Body())).
					WithField("status", response.StatusCode()).
					Info("sent request")
			}
		}
	})

	g.Go(func() error {
		gr, _ := errgroup.WithContext(ctx)

		for {
			select {
			case <-gCtx.Done():
				log.Info("stopped send metrics")
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

							log.WithField("uri", response.Request.URL).
								WithField("response", string(response.Body())).
								WithField("status", response.StatusCode()).
								Info("sent request")
						}

						return nil
					})
				}

				if err := gr.Wait(); err != nil {
					log.Error(err.Error())
				}
			}
		}
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
