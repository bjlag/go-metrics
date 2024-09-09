//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package counter

type Storage interface {
	AddCounter(name string, value int64)
}
