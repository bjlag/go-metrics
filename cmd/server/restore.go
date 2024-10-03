package main

import (
	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage"
	"github.com/bjlag/go-metrics/internal/storage/file"
)

func restoreData(fileStorage *file.Storage, memStorage storage.Repository) error {
	data, err := fileStorage.Load()
	if err != nil {
		return err
	}

	for _, value := range data {
		switch value.MType {
		case model.TypeCounter:
			memStorage.AddCounter(value.ID, *value.Delta)
		case model.TypeGauge:
			memStorage.SetGauge(value.ID, *value.Value)
		}
	}

	return nil
}
