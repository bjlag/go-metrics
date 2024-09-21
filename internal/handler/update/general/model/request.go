package model

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

type Request struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m Request) IsValid() bool {
	if m.IsGauge() && m.Value != nil {
		return true
	}

	if m.IsCounter() && m.Delta != nil {
		return true
	}

	return false
}

func (m Request) IsGauge() bool {
	return m.MType == TypeGauge
}

func (m Request) IsCounter() bool {
	return m.MType == TypeCounter
}
