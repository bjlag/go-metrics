package memory

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

func (s *MemStorage) GetGauge(name string) float64 {
	return s.Gauge[name]
}

func (s *MemStorage) SetGauge(name string, value float64) {
	s.Gauge[name] = value
}

func (s *MemStorage) GetCounter(name string) int64 {
	return s.Counter[name]
}

func (s *MemStorage) AddCounter(name string, value int64) {
	currentValue, ok := s.Counter[name]
	if !ok {
		s.Counter[name] = value
		return
	}

	s.Counter[name] = currentValue + value
}
