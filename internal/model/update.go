package model

import (
	"encoding/json"
	"errors"
)

// UpdateIn модель описывает входящий запрос на обновление метрики.
type UpdateIn struct {
	ID    string   `json:"id" example:"Sys"`              // Имя метрики
	MType string   `json:"type" example:"gauge"`          // Параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty" example:"1"`   // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty" example:"1.1"` // Значение метрики в случае передачи gauge
}

// IsValid проверяет валидный ли запрос.
func (m *UpdateIn) IsValid() bool {
	if m.IsGauge() && m.Value != nil {
		return true
	}

	if m.IsCounter() && m.Delta != nil {
		return true
	}

	return false
}

// IsGauge возвращает true, если запрос с типом метрики [TypeGauge]
func (m *UpdateIn) IsGauge() bool {
	return m.MType == TypeGauge
}

// IsCounter возвращает true, если запрос с типом метрики [TypeCounter]
func (m *UpdateIn) IsCounter() bool {
	return m.MType == TypeCounter
}

// UnmarshalJSON анмаршалинг запроса в модель [UpdateIn] с валидацией входящих данных.
func (m *UpdateIn) UnmarshalJSON(b []byte) error {
	type UpdateInAlias UpdateIn

	aliasValue := &struct {
		*UpdateInAlias
	}{
		UpdateInAlias: (*UpdateInAlias)(m),
	}

	err := json.Unmarshal(b, &aliasValue)
	if err != nil {
		return err
	}

	if m.ID == "" {
		return ErrInvalidID
	}

	var errs []error
	if !m.IsValid() {
		errs = append(errs, ErrInvalidID)
	}

	if m.IsCounter() && m.Delta == nil {
		errs = append(errs, ErrInvalidValue)
	}

	if m.IsGauge() && m.Value == nil {
		errs = append(errs, ErrInvalidValue)
	}

	return errors.Join(errs...)
}

// UpdateOut модель описывает ответ результата обновления метрики.
type UpdateOut struct {
	ID    string   `json:"id"`              // Имя метрики
	MType string   `json:"type"`            // Параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // Значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // Значение метрики в случае передачи gauge
}
