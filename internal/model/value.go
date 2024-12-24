package model

import (
	"encoding/json"
	"errors"
)

// ValueIn модель описывает входящий запрос на получение значения метрики.
type ValueIn struct {
	ID    string `json:"id"`   // имя метрики
	MType string `json:"type"` // параметр, принимающий значение gauge или counter
}

// IsValid проверяет валидный ли запрос.
func (m *ValueIn) IsValid() bool {
	if m.IsGauge() {
		return true
	}

	if m.IsCounter() {
		return true
	}

	return false
}

// IsGauge возвращает true, если запрос с типом метрики [TypeGauge]
func (m *ValueIn) IsGauge() bool {
	return m.MType == TypeGauge
}

// IsCounter возвращает true, если запрос с типом метрики [TypeCounter]
func (m *ValueIn) IsCounter() bool {
	return m.MType == TypeCounter
}

// UnmarshalJSON анмаршалинг запроса в модель [ValueIn] с валидацией входящих данных.
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

// ValueOut модель описывает ответ результата получения значения метрики.
type ValueOut struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
