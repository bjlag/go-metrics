package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type Storage struct {
	lock sync.RWMutex
	path string
}

func NewStorage(path string) (*Storage, error) {
	if len(path) == 0 {
		return nil, errors.New("path cannot be empty")
	}

	return &Storage{
		path: path,
	}, nil
}

func (s *Storage) Save(data []Metric) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling data: %w", err)
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	return os.WriteFile(s.path, bytes, 0666)
}

func (s *Storage) Load() ([]Metric, error) {
	if _, err := os.Stat(s.path); os.IsNotExist(err) {
		return nil, nil
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	if len(data) == 0 {
		return nil, nil
	}

	var metrics []Metric
	err = json.Unmarshal(data, &metrics)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling file: %w", err)
	}

	return metrics, nil
}
