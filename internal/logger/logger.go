package logger

type Logger interface {
	WithField(key string, value interface{}) Logger
	WithError(err error) Logger
	Error(msg string)
	Info(msg string)
	Debug(msg string)
}
