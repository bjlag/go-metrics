package memory

import (
	"context"
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

func (s *Storage) GetAllGauges(_ context.Context) storage.Gauges {
	return s.gauges
}

func (s *Storage) GetAllCounters(_ context.Context) storage.Counters {
	return s.counters
}

func (s *Storage) GetGauge(_ context.Context, id string) (float64, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, ok := s.gauges[id]
	if !ok {
		return 0, storage.NewMetricNotFoundError("gauge", id)
	}

	return value, nil
}

func (s *Storage) SetGauge(_ context.Context, id string, value float64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.gauges[id] = value
}

func (s *Storage) GetCounter(_ context.Context, id string) (int64, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, ok := s.counters[id]
	if !ok {
		return 0, storage.NewMetricNotFoundError("counter", id)
	}

	return value, nil
}

func (s *Storage) AddCounter(_ context.Context, id string, value int64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	currentValue, ok := s.counters[id]
	if !ok {
		s.counters[id] = value
		return
	}

	s.counters[id] = currentValue + value
}
