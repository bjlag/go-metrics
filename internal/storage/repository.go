package storage

import "context"

// Gauges тип для метрик типа gauge.
type Gauges map[string]float64

// Counters тип для метрик типа counter.
type Counters map[string]int64

// Gauge тип метрики gauge.
type Gauge struct {
	// ID метрики.
	ID string
	// Value значение метрики.
	Value float64
}

// Counter тип метрики counter.
type Counter struct {
	// ID метрики.
	ID string
	// Value значение метрики.
	Value int64
}

// Repository общий интерфейс репозитория для работы с метриками.
type Repository interface {
	// GetAllGauges возвращает все метрики типа Gauge, которые хранятся в хранилище.
	GetAllGauges(ctx context.Context) Gauges
	// GetAllCounters возвращает все метрики типа Counter, которые хранятся в хранилище.
	GetAllCounters(ctx context.Context) Counters
	// GetGauge возвращает значение метрики типа Gauge по ее ID.
	GetGauge(ctx context.Context, id string) (float64, error)
	// SetGauge записывает переданное значение метрики типа Gauge по ее ID в хранилище.
	SetGauge(ctx context.Context, id string, value float64)
	// SetGauges записывает набор переданных метрик типа Gauge в хранилище.
	SetGauges(ctx context.Context, gauges []Gauge) error
	// GetCounter возвращает значение метрики типа Counter по ее ID.
	GetCounter(ctx context.Context, id string) (int64, error)
	// AddCounter добавляет переданное значение метрики типа Counter по ее ID в хранилище.
	AddCounter(ctx context.Context, id string, value int64)
	// AddCounters добавляет значения из набора переданных метрик типа Counter в хранилище.
	AddCounters(ctx context.Context, counters []Counter) error
}
