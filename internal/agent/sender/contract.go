//go:generate mockgen -source ${GOFILE} -package mock -destination mock/contract_mock.go

package sender

import "net/http"

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}
