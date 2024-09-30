package async

import (
	"time"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage/file"
	"github.com/bjlag/go-metrics/internal/storage/memory"
)

type Backup struct {
	storage  *memory.Storage
	fStorage *file.Storage
	interval time.Duration
	log      logger.Logger

	ticker     *time.Ticker
	needUpdate bool
}

func New(storage *memory.Storage, fStorage *file.Storage, interval time.Duration, log logger.Logger) *Backup {
	return &Backup{
		storage:  storage,
		fStorage: fStorage,
		interval: interval,
		log:      log,
	}
}

func (b *Backup) Start() {
	b.ticker = time.NewTicker(b.interval)

	go func() {
		for range b.ticker.C {
			if b.needUpdate {
				err := b.update()
				if err != nil {
					b.log.WithField("error", err.Error()).
						Error("failed to update backup")
				}

				b.needUpdate = false
			}
		}
	}()

	b.log.Info("async backup started")
}

func (b *Backup) Stop() {
	b.ticker.Stop()

	err := b.update()
	if err != nil {
		b.log.WithField("error", err.Error()).
			Error("failed to update backup while stopping")
	}

	b.log.Info("backup stopped")
}

func (b *Backup) Create() error {
	b.needUpdate = true

	return nil
}

func (b *Backup) update() error {
	counters := b.storage.GetAllCounters()
	gauges := b.storage.GetAllGauges()

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
		b.log.WithField("error", err.Error()).
			Error("failed to backup data")
		return err
	}

	return nil
}
