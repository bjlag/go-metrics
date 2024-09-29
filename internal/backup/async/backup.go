package sync

import (
	"fmt"
	"sync"
	"time"

	"github.com/bjlag/go-metrics/internal/logger"
	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage/file"
	"github.com/bjlag/go-metrics/internal/storage/memory"
)

type Backup struct {
	lock     sync.RWMutex
	storage  *memory.Storage
	fStorage *file.Storage
	interval time.Duration
	log      logger.Logger
}

func New(storage *memory.Storage, fStorage *file.Storage, interval time.Duration, log logger.Logger) *Backup {
	return &Backup{
		storage:  storage,
		fStorage: fStorage,
		interval: interval,
		log:      log,
	}
}

func (b *Backup) Create() error {
	if b.interval == 0 {
		err := b.update()
		if err != nil {
			return err
		}

		return nil
	}

	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			err := b.update()
			if err != nil {

			}
		}
	}

	return nil
}

func (b *Backup) update() error  {
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
