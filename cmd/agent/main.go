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
	"github.com/bjlag/go-metrics/internal/agent/limiter"
	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/signature"
)

func main() {
	parseFlags()
	parseEnvs()

	log, err := logger.NewZapLog(logLevel)
	if err != nil {
		logNativ.Fatalln(err)
	}
	defer func() {
		_ = log.Close()
	}()

	log.Info("Starting agent")
	log.Info(fmt.Sprintf("Sending metrics to %s", addr.String()))
	log.Info(fmt.Sprintf("Poll interval is %s", pollInterval))
	log.Info(fmt.Sprintf("Report interval is %s", reportInterval))
	log.Info(fmt.Sprintf("Log level is '%s'", logLevel))
	log.Info(fmt.Sprintf("Sign request is %t", len(secretKey) > 0))
	log.Info(fmt.Sprintf("Rate limit is %d", rateLimit))

	if err := run(log); err != nil {
		log.WithError(err).Error("Error running agent")
	}
}

func run(log logger.Logger) error {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	signManager := signature.NewSignManager(secretKey)
	rateLimiter := limiter.NewRateLimiter(rateLimit)
	metricCollector := collector.NewMetricCollector(&runtime.MemStats{})
	metricClient := client.NewHTTPSender(addr.host, addr.port, signManager, rateLimiter, log)

	pollTicker := time.NewTicker(pollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(reportInterval)
	defer reportTicker.Stop()

	g, gCtx := errgroup.WithContext(ctx)

	go func() {
		<-ctx.Done()

		log.Info("Graceful shutting down agent")
	}()

	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				log.Info("Stopped read stats")
				return nil
			case <-pollTicker.C:
				metricCollector.ReadStats()

				metrics := []*collector.Metric{
					collector.NewCounterMetric("PollCount", 1),
				}

				if err := metricClient.Send(metrics); err != nil {
					log.WithError(err).Error("Error in sending poll count")
				}
			}
		}
	})

	g.Go(func() error {
		for {
			select {
			case <-gCtx.Done():
				log.Info("Stopped send metrics")
				return nil
			case <-reportTicker.C:
				metrics := metricCollector.Collect()
				if len(metrics) == 0 {
					continue
				}

				if err := metricClient.Send(metrics); err != nil {
					log.WithError(err).Error("Error in sending report")
				}
			}
		}
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
