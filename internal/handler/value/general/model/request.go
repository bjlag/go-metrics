package model

const (
	TypeGauge   = "gauge"
	TypeCounter = "counter"
)

type Request struct {
	ID    string `json:"id"`   // имя метрики
	MType string `json:"type"` // параметр, принимающий значение gauge или counter
}

func (m Request) IsValid() bool {
	if m.IsGauge() {
		return true
	}

	if m.IsCounter() {
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
