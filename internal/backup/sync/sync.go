package sync

import (
	"context"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage"
	"github.com/bjlag/go-metrics/internal/storage/file"
)

// Backup обслуживает создание синхронной резервной копии метрик.
type Backup struct {
	storage  storage.Repository
	fStorage *file.Storage
	log      logger.Logger
}

// New создает экземпляр сервиса по созданию резервных копий.
func New(storage storage.Repository, fStorage *file.Storage, log logger.Logger) *Backup {
	return &Backup{
		storage:  storage,
		fStorage: fStorage,
		log:      log,
	}
}

// Create создает резервную копию.
func (b *Backup) Create(ctx context.Context) error {
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
