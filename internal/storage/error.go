package storage

import "fmt"

// NotFoundError описывает ошибку когда метрика не найдена.
type NotFoundError struct {
	kind string
	name string
	err  error
}

// NewMetricNotFoundError создает ошибку.
func NewMetricNotFoundError(kind string, name string, err error) *NotFoundError {
	return &NotFoundError{
		kind: kind,
		name: name,
		err:  err,
	}
}

// Error возвращает текст ошибки.
func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s metric '%s' not found", e.kind, e.name)
}

// Unwrap возвращает следующую ошибку в цепочке ошибок.
func (e NotFoundError) Unwrap() error {
	return e.err
}

// Kind возвращает тип метрики.
func (e NotFoundError) Kind() string {
	return e.kind
}

// Name возвращает название метрики.
func (e NotFoundError) Name() string {
	return e.name
}
