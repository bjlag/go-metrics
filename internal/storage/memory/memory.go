package memory

import (
	"fmt"

	"github.com/bjlag/go-metrics/internal/storage"
)

const (
	initSize = 10
)

type MemStorage struct {
	gauges   storage.Gauges
	counters storage.Counters
}

func NewMemStorage() *MemStorage {
	gauges := make(storage.Gauges, initSize)
	counters := make(storage.Counters, initSize)

	return &MemStorage{
		gauges:   gauges,
		counters: counters,
	}
}

func (s *MemStorage) GetAllGauges() storage.Gauges {
	return s.gauges
}

func (s *MemStorage) GetAllCounters() storage.Counters {
	return s.counters
}

func (s *MemStorage) GetGauge(name string) (float64, error) {
	value, ok := s.gauges[name]
	if !ok {
		return 0, fmt.Errorf("gauge metric '%s' not found", name)
	}

	return value, nil
}

func (s *MemStorage) SetGauge(name string, value float64) {
	s.gauges[name] = value
}

func (s *MemStorage) GetCounter(name string) (int64, error) {
	value, ok := s.counters[name]
	if !ok {
		return 0, fmt.Errorf("counter metric '%s' not found", name)
	}

	return value, nil
}

func (s *MemStorage) AddCounter(name string, value int64) {
	currentValue, ok := s.counters[name]
	if !ok {
		s.counters[name] = value
		return
	}

	s.counters[name] = currentValue + value
}
