package model

import "errors"

var (
	// ErrInvalidID ошибка, если в запросе не указан ID метрики.
	ErrInvalidID = errors.New("metric ID not specified")
	// ErrInvalidType ошибка, если в запросе передан невалидный тип метрики.
	ErrInvalidType = errors.New("metric type is invalid")
	// ErrInvalidValue ошибка, если в запросе передано невалидное значение.
	ErrInvalidValue = errors.New("metric value is invalid")
)
