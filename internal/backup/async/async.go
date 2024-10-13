package async

import (
	"context"
	"time"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage"
	"github.com/bjlag/go-metrics/internal/storage/file"
)

type Backup struct {
	storage  storage.Repository
	fStorage *file.Storage
	interval time.Duration
	log      logger.Logger

	ticker     *time.Ticker
	needUpdate bool
}

func New(storage storage.Repository, fStorage *file.Storage, interval time.Duration, log logger.Logger) *Backup {
	return &Backup{
		storage:  storage,
		fStorage: fStorage,
		interval: interval,
		log:      log,
	}
}

func (b *Backup) Start(ctx context.Context) {
	b.ticker = time.NewTicker(b.interval)

	go func() {
		for range b.ticker.C {
			if b.needUpdate {
				err := b.update(ctx)
				if err != nil {
					b.log.WithError(err).Error("Failed to update backup")
				}

				b.needUpdate = false
			}
		}
	}()

	b.log.Info("async backup started")
}

func (b *Backup) Stop(ctx context.Context) {
	b.ticker.Stop()

	err := b.update(ctx)
	if err != nil {
		b.log.WithError(err).Error("Failed to update backup while stopping")
	}

	b.log.Info("Backup stopped")
}

func (b *Backup) Create(_ context.Context) error {
	b.needUpdate = true

	return nil
}

func (b *Backup) update(ctx context.Context) error {
	counters := b.storage.GetAllCounters(ctx)
	gauges := b.storage.GetAllGauges(ctx)

	data := make([]file.Metric, 0, len(counters)+len(gauges))

	for id, value := range counters {
		data = append(data, file.Metric{
			ID:    id,
			MType: model.TypeCounter,
			Delta: &value,
		})
	}

	for id, value := range gauges {
		data = append(data, file.Metric{
			ID:    id,
			MType: model.TypeGauge,
			Value: &value,
		})
	}

	err := b.fStorage.Save(data)
	if err != nil {
		b.log.WithError(err).Error("Failed to backup data")
		return err
	}

	return nil
}
