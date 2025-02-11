package main

import (
	"context"
	"fmt"
	"github.com/bjlag/go-metrics/internal/agent/client/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	logNativ "log"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/bjlag/go-metrics/cmd"
	"github.com/bjlag/go-metrics/cmd/agent/config"
	"github.com/bjlag/go-metrics/internal/agent/collector"
	"github.com/bjlag/go-metrics/internal/logger"
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
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	//cryptManager, err := crypt.NewEncryptManager(cfg.CryptoKeyPath)
	//if err != nil {
	//	return err
	//}

	//signManager := signature.NewSignManager(cfg.SecretKey)
	//rateLimiter := limiter.NewRateLimiter(cfg.RateLimit)
	metricCollector := collector.NewMetricCollector(&runtime.MemStats{})
	//metricClient := http.NewSender(cfg.Address.Host, cfg.Address.Port, signManager, cryptManager, rateLimiter, log)

	grpcConn, err := grpc.NewClient(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer func() {
		_ = grpcConn.Close()
	}()

	metricClient := rpc.NewSender(grpcConn)

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

	log.Info("Agent shutdown gracefully")

	return nil
}
