package memory

import (
	"context"
	"sync"

	"github.com/bjlag/go-metrics/internal/storage"
)

const (
	initSize = 100
)

// Storage обслуживает in-memory хранилище.
type Storage struct {
	lock     sync.RWMutex
	gauges   storage.Gauges
	counters storage.Counters
}

// NewStorage создает хранилище.
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
		return 0, storage.NewMetricNotFoundError("gauge", id, nil)
	}

	return value, nil
}

func (s *Storage) SetGauge(_ context.Context, id string, value float64) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.gauges[id] = value
}

func (s *Storage) SetGauges(_ context.Context, gauges []storage.Gauge) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, gauge := range gauges {
		s.gauges[gauge.ID] = gauge.Value
	}

	return nil
}

func (s *Storage) GetCounter(_ context.Context, id string) (int64, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()

	value, ok := s.counters[id]
	if !ok {
		return 0, storage.NewMetricNotFoundError("counter", id, nil)
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

func (s *Storage) AddCounters(_ context.Context, counters []storage.Counter) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	for _, counter := range counters {
		currentValue, ok := s.counters[counter.ID]
		if !ok {
			s.counters[counter.ID] = counter.Value
			continue
		}

		s.counters[counter.ID] = currentValue + counter.Value
	}

	return nil
}
