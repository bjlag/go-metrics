package storage

type Repository interface {
	GetGauge(name string) float64
	SetGauge(name string, value float64)
	GetCounter(name string) int64
	AddCounter(name string, value int64)
}
