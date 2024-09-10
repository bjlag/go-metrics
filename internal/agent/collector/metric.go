package collector

const (
	Gauge   = "gauge"
	Counter = "counter"
)

type Metric struct {
	kind  string
	name  string
	value interface{}
}

func NewMetric(kind string, name string, value interface{}) *Metric {
	return &Metric{
		kind:  kind,
		name:  name,
		value: value,
	}
}

func (m Metric) Kind() string {
	return m.kind
}

func (m Metric) Name() string {
	return m.name
}

func (m Metric) Value() interface{} {
	return m.value
}
