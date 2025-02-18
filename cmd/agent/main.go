package main

import (
	"context"
	"fmt"
	logNativ "log"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/bjlag/go-metrics/cmd"
	"github.com/bjlag/go-metrics/cmd/agent/config"
	agent "github.com/bjlag/go-metrics/internal/agent/client"
	"github.com/bjlag/go-metrics/internal/agent/client/http"
	"github.com/bjlag/go-metrics/internal/agent/client/rpc"
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

	protocol := "HTTP"
	address := cfg.AddressHTTP
	if cfg.AddressRPC != nil {
		protocol = "RPC"
		address = cfg.AddressRPC
	}

	log.Info(fmt.Sprintf("Sending metrics to %s server: %s", protocol, address.String()))
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
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	cryptManager, err := crypt.NewEncryptManager(cfg.CryptoKeyPath)
	if err != nil {
		return err
	}

	signManager := signature.NewSignManager(cfg.SecretKey)
	rateLimiter := limiter.NewRateLimiter(cfg.RateLimit)
	metricCollector := collector.NewMetricCollector(&runtime.MemStats{})

	var client agent.Client

	if cfg.AddressRPC != nil {
		client = rpc.NewSender(cfg.AddressRPC.String(), signManager, log)
		defer func() {
			_ = client.(*rpc.MetricSender).Close()
		}()
	}

	if client == nil && cfg.AddressHTTP != nil {
		client = http.NewSender(cfg.AddressHTTP.Host, cfg.AddressHTTP.Port, signManager, cryptManager, rateLimiter, log)
	}

	if client == nil {
		return fmt.Errorf("could not create client")
	}

	pollTicker := time.NewTicker(cfg.PollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(cfg.ReportInterval)
	defer reportTicker.Stop()

	g, gCtx := errgroup.WithContext(ctx)

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

				if err := client.Send(metrics); err != nil {
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

				if err := client.Send(metrics); err != nil {
					log.WithError(err).Error("Error in sending report")
				}
			}
		}
	})

	if err := g.Wait(); err != nil {
		return err
	}

	log.Info("Agent shutdown gracefully")

	return nil
}
