package async

import (
	"context"
	"time"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage"
	"github.com/bjlag/go-metrics/internal/storage/file"
)

// Backup обслуживает создание асинхронной резервной копии метрик.
type Backup struct {
	storage  storage.Repository
	fStorage *file.Storage
	interval time.Duration
	log      logger.Logger

	ticker     *time.Ticker
	needUpdate bool
}

// New создает экземпляр сервиса по созданию резервных копий.
// Параметр interval регулирует, с какой периодичностью надо делать резервную копию.
func New(storage storage.Repository, fStorage *file.Storage, interval time.Duration, log logger.Logger) *Backup {
	return &Backup{
		storage:  storage,
		fStorage: fStorage,
		interval: interval,
		log:      log,
	}
}

// Start запускает воркер, которая в фоновом режиме создает резервные копии.
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

// Stop останавливает асинхронный воркер, создающий резервные копии.
func (b *Backup) Stop(ctx context.Context) {
	b.ticker.Stop()

	err := b.update(ctx)
	if err != nil {
		b.log.WithError(err).Error("Failed to update backup while stopping")
	}

	b.log.Info("Backup stopped")
}

// Create посылает сигнал, что надо создать копию данных.
func (b *Backup) Create(_ context.Context) error {
	b.needUpdate = true

	return nil
}

// Функция update обновляет резервную копию.
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
