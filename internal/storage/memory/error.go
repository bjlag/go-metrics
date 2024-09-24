package memory

import "fmt"

type MetricNotFoundError struct {
	kind string
	name string
}

func NewMetricNotFoundError(kind string, name string) *MetricNotFoundError {
	return &MetricNotFoundError{
		kind: kind,
		name: name,
	}
}

func (e MetricNotFoundError) Error() string {
	return fmt.Sprintf("%s metric '%s' not found", e.kind, e.name)
}
