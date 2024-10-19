package collector

import (
	"fmt"
	"strconv"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Metric struct {
	mType string
	name  string
	value interface{}
}

func NewMetric(mType string, name string, value interface{}) *Metric {
	return &Metric{
		mType: mType,
		name:  name,
		value: value,
	}
}

func NewCounterMetric(name string, value interface{}) *Metric {
	return NewMetric(Counter, name, value)
}

func NewGaugeMetric(name string, value interface{}) *Metric {
	return NewMetric(Gauge, name, value)
}

func (m Metric) Kind() string {
	return m.mType
}

func (m Metric) Name() string {
	return m.name
}

func (m Metric) Value() interface{} {
	return m.value
}

func (m Metric) GaugeValue() (float64, error) {
	switch v := m.value.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case string:
		if value, err := strconv.ParseFloat(v, 64); err == nil {
			return value, nil
		}
		return 0, fmt.Errorf("failed to convert value to float64 from string: %s", v)
	default:
		return 0, fmt.Errorf("unknow type value: %s", v)
	}
}

func (m Metric) CounterValue() (int64, error) {
	switch v := m.value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case string:
		if value, err := strconv.ParseInt(v, 10, 64); err == nil {
			return value, nil
		}
		return 0, fmt.Errorf("failed to convert value to int64 from string: %s", v)
	default:
		return 0, fmt.Errorf("unknow type value: %s", v)
	}
}
