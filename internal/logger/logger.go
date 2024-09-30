package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	WithField(key string, value interface{}) Logger
	Error(msg string)
	Info(msg string)
	Debug(msg string)
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

func (l *ZapLogger) WithField(key string, value interface{}) Logger {
	return &ZapLogger{
		logger: l.logger.With(zap.Any(key, value)),
	}
}

func (l *ZapLogger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *ZapLogger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *ZapLogger) Debug(msg string) {
	l.logger.Debug(msg)
}
