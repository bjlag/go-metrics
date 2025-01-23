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

	"github.com/bjlag/go-metrics/cmd"
	"github.com/bjlag/go-metrics/cmd/agent/config"
	"github.com/bjlag/go-metrics/internal/agent/client"
	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/agent/limiter"
	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/securety/crypt"
	"github.com/bjlag/go-metrics/internal/securety/signature"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	cfg := config.LoadConfig()

	log, err := logger.NewZapLog(cfg.LogLevel)
	if err != nil {
		logNativ.Fatalln(err)
	}
	defer func() {
		_ = log.Close()
	}()

	build := cmd.NewBuild(buildVersion, buildDate, buildCommit)
	log.Info(build.VersionString())
	log.Info(build.DateString())
	log.Info(build.CommitString())

	log.Info("Starting agent")
	log.Info(fmt.Sprintf("Sending metrics to %s", cfg.Address.String()))
	log.Info(fmt.Sprintf("Poll interval is %s", cfg.PollInterval))
	log.Info(fmt.Sprintf("Report interval is %s", cfg.ReportInterval))
	log.Info(fmt.Sprintf("Log level is '%s'", cfg.LogLevel))
	log.Info(fmt.Sprintf("Sign request is %t", len(cfg.SecretKey) > 0))
	log.Info(fmt.Sprintf("Rate limit is %d", cfg.RateLimit))
	log.Info(fmt.Sprintf("Public key %s", cfg.CryptoKeyPath))
	log.Info(fmt.Sprintf("JSON config %s", cfg.ConfigPath))

	if err := run(log, cfg); err != nil {
		log.WithError(err).Error("Error running agent")
	}
}

func run(log logger.Logger, cfg *config.Configuration) error {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	cryptManager, err := crypt.NewEncryptManager(cfg.CryptoKeyPath)
	if err != nil {
		return err
	}

	signManager := signature.NewSignManager(cfg.SecretKey)
	rateLimiter := limiter.NewRateLimiter(cfg.RateLimit)
	metricCollector := collector.NewMetricCollector(&runtime.MemStats{})
	metricClient := client.NewHTTPSender(cfg.Address.Host, cfg.Address.Port, signManager, cryptManager, rateLimiter, log)

	pollTicker := time.NewTicker(cfg.PollInterval)
	defer pollTicker.Stop()
	reportTicker := time.NewTicker(cfg.ReportInterval)
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
				metrics, err := metricCollector.Collect()
				if err != nil {
					log.WithError(err).Error("Error in getting metrics")
				}

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
