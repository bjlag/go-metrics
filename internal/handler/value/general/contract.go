//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package general

type Storage interface {
	GetGauge(name string) (float64, error)
	GetCounter(name string) (int64, error)
}