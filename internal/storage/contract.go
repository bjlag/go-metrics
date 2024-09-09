package storage

type Repository interface {
	GetGauge(name string) (float64, error)
	SetGauge(name string, value float64)
	GetCounter(name string) (int64, error)
	AddCounter(name string, value int64)
}
