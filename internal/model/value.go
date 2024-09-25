package model

type ValueIn struct {
	ID    string `json:"id"`   // имя метрики
	MType string `json:"type"` // параметр, принимающий значение gauge или counter
}

func (m ValueIn) IsValid() bool {
	if m.IsGauge() {
		return true
	}

	if m.IsCounter() {
		return true
	}

	return false
}

func (m ValueIn) IsGauge() bool {
	return m.MType == TypeGauge
}

func (m ValueIn) IsCounter() bool {
	return m.MType == TypeCounter
}

type ValueOut struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}
