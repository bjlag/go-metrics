package storage

import "context"

type Gauges map[string]float64
type Counters map[string]int64

type Repository interface {
	GetAllGauges(ctx context.Context) Gauges
	GetAllCounters(ctx context.Context) Counters
	GetGauge(ctx context.Context, id string) (float64, error)
	SetGauge(ctx context.Context, id string, value float64)
	GetCounter(ctx context.Context, id string) (int64, error)
	AddCounter(ctx context.Context, id string, value int64)
}
