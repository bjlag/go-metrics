package storage

type Interface interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)
}
