package memory

import (
	"fmt"
)

const (
	initSize = 10
)

type Gauge map[string]float64
type Counter map[string]int64

type MemStorage struct {
	Gauge   Gauge
	Counter Counter
}

func NewMemStorage() *MemStorage {
	gauge := make(Gauge, initSize)
	counter := make(Counter, initSize)

	return &MemStorage{
		Gauge:   gauge,
		Counter: counter,
	}
}

func (s *MemStorage) GetGauge(name string) (float64, error) {
	value, ok := s.Gauge[name]
	if !ok {
		return 0, fmt.Errorf("gauge metric '%s' not found", name)
	}

	return value, nil
}

func (s *MemStorage) SetGauge(name string, value float64) {
	s.Gauge[name] = value
}

func (s *MemStorage) GetCounter(name string) (int64, error) {
	value, ok := s.Counter[name]
	if !ok {
		return 0, fmt.Errorf("counter metric '%s' not found", name)
	}

	return value, nil
}

func (s *MemStorage) AddCounter(name string, value int64) {
	currentValue, ok := s.Counter[name]
	if !ok {
		s.Counter[name] = value
		return
	}

	s.Counter[name] = currentValue + value
}
