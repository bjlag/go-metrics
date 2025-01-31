//go:generate mockgen -source ${GOFILE} -package mock -destination ../mock/logger_mock.go

package logger

// Logger интерфейс логгера.
type Logger interface {
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	Error(msg string)
	Info(msg string)
	Debug(msg string)
}
