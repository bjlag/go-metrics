package storage

import "fmt"

type NotFoundError struct {
	kind string
	name string
	err  error
}

func NewMetricNotFoundError(kind string, name string, err error) *NotFoundError {
	return &NotFoundError{
		kind: kind,
		name: name,
		err:  err,
	}
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s metric '%s' not found", e.kind, e.name)
}

func (e NotFoundError) Unwrap() error {
	return e.err
}

func (e NotFoundError) Kind() string {
	return e.kind
}

func (e NotFoundError) Name() string {
	return e.name
}
