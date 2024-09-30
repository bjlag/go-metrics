package storage

type Gauges map[string]float64
type Counters map[string]int64

type Repository interface {
	GetAllGauges() Gauges
	GetAllCounters() Counters
	GetGauge(name string) (float64, error)
	SetGauge(name string, value float64)
	GetCounter(name string) (int64, error)
	AddCounter(name string, value int64)
}
