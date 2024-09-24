//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package unknown

type Logger interface {
	Info(msg string, fields map[string]interface{})
}
