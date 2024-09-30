//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package counter

type Storage interface {
	AddCounter(name string, value int64)
}

type Backup interface {
	Create() error
}

type Logger interface {
	Error(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
}
