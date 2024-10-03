package model

type UpdateIn struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func (m UpdateIn) IsValid() bool {
	if m.IsGauge() && m.Value != nil {
		return true
	}

	if m.IsCounter() && m.Delta != nil {
		return true
	}

	return false
}

func (m UpdateIn) IsGauge() bool {
	return m.MType == TypeGauge
}

func (m UpdateIn) IsCounter() bool {
	return m.MType == TypeCounter
}

type UpdateOut struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
