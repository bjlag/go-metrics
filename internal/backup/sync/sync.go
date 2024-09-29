package sync

import (
	"fmt"
	"sync"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage/file"
	"github.com/bjlag/go-metrics/internal/storage/memory"
)

type Backup struct {
	lock     sync.RWMutex
	storage  *memory.Storage
	fStorage *file.Storage
	log      logger.Logger
}

func New(storage *memory.Storage, fStorage *file.Storage, log logger.Logger) *Backup {
	return &Backup{
		storage:  storage,
		fStorage: fStorage,
		log:      log,
	}
}

func (b *Backup) Create() error {
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
		b.log.Error(fmt.Sprintf("Failed to backup data: %s", err.Error()), nil)
		return err
	}

	return nil
}
