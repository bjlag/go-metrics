package main

import (
	"context"

	"github.com/bjlag/go-metrics/internal/model"
	"github.com/bjlag/go-metrics/internal/storage"
	"github.com/bjlag/go-metrics/internal/storage/file"
)

func restoreData(ctx context.Context, fileStorage *file.Storage, memStorage storage.Repository) error {
	data, err := fileStorage.Load()
	if err != nil {
		return err
	}

	for _, value := range data {
		switch value.MType {
		case model.TypeCounter:
			memStorage.AddCounter(ctx, value.ID, *value.Delta)
		case model.TypeGauge:
			memStorage.SetGauge(ctx, value.ID, *value.Value)
		}
	}

	return nil
}
