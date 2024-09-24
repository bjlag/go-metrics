//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package general

type Storage interface {
	SetGauge(name string, value float64)
	AddCounter(name string, value int64)
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}

type Logger interface {
	Error(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
}
