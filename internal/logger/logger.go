package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Error(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Debug(msg string, fields map[string]interface{})
}

type ZapLogger struct {
	logger *zap.Logger
}

func NewZapLogger(level string) (*ZapLogger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	cfg.DisableCaller = true
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLogger{
		logger: logger,
	}, nil
}

func (l *ZapLogger) Close() error {
	if err := l.logger.Sync(); err != nil {
		return err
	}

	return nil
}

func (l *ZapLogger) Error(msg string, fields map[string]interface{}) {
	l.logger.Error(msg, getZapFields(fields)...)
}

func (l *ZapLogger) Info(msg string, fields map[string]interface{}) {
	l.logger.Info(msg, getZapFields(fields)...)
}

func (l *ZapLogger) Debug(msg string, fields map[string]interface{}) {
	l.logger.Debug(msg, getZapFields(fields)...)
}

func getZapFields(fields map[string]interface{}) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for filed, value := range fields {
		zapFields = append(zapFields, zap.Any(filed, value))
	}

	return zapFields
}
