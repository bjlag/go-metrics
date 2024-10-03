package memory

import (
	"sync"

	"github.com/bjlag/go-metrics/internal/storage"
)

const (
	initSize = 100
)

type Storage struct {
	lock     sync.RWMutex
	gauges   storage.Gauges
	counters storage.Counters
}

func NewStorage() *Storage {
	gauges := make(storage.Gauges, initSize)
	counters := make(storage.Counters, initSize)

	return &Storage{
		gauges:   gauges,
		counters: counters,
	}
}

func (s *Storage) GetAllGauges() storage.Gauges {
	return s.gauges
}

func (s *Storage) GetAllCounters() storage.Counters {
	return s.counters
}

func (s *Storage) GetGauge(name string) (float64, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, ok := s.gauges[name]
	if !ok {
		return 0, NewMetricNotFoundError("gauge", name)
	}

	return value, nil
}

func (s *Storage) SetGauge(name string, value float64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.gauges[name] = value
}

func (s *Storage) GetCounter(name string) (int64, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, ok := s.counters[name]
	if !ok {
		return 0, NewMetricNotFoundError("counter", name)
	}

	return value, nil
}

func (s *Storage) AddCounter(name string, value int64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	currentValue, ok := s.counters[name]
	if !ok {
		s.counters[name] = value
		return
	}

	s.counters[name] = currentValue + value
}
