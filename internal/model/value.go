package model

import (
	"encoding/json"
	"errors"
)

type ValueIn struct {
	ID    string `json:"id"`   // имя метрики
	MType string `json:"type"` // параметр, принимающий значение gauge или counter
}

func (m *ValueIn) IsValid() bool {
	if m.IsGauge() {
		return true
	}

	if m.IsCounter() {
		return true
	}

	return false
}

func (m *ValueIn) IsGauge() bool {
	return m.MType == TypeGauge
}

func (m *ValueIn) IsCounter() bool {
	return m.MType == TypeCounter
}

func (m *ValueIn) UnmarshalJSON(b []byte) error {
	type ValueInAlias ValueIn

	aliasValue := &struct {
		*ValueInAlias
	}{
		ValueInAlias: (*ValueInAlias)(m),
	}

	err := json.Unmarshal(b, &aliasValue)
	if err != nil {
		return err
	}

	var errs []error
	if m.ID == "" {
		errs = append(errs, ErrInvalidID)
	}

	if !m.IsValid() {
		errs = append(errs, ErrInvalidType)
	}

	return errors.Join(errs...)
}

type ValueOut struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
