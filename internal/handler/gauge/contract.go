//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package gauge

type Storage interface {
	SetGauge(name string, value float64)
}
