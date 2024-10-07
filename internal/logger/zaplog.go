package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLog struct {
	logger *zap.Logger
}

func NewZapLog(level string) (*ZapLog, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	cfg.DisableCaller = true
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	return &ZapLog{
		logger: logger,
	}, nil
}

func (l *ZapLog) Close() error {
	if err := l.logger.Sync(); err != nil {
		return err
	}

	return nil
}

func (l *ZapLog) WithField(key string, value interface{}) Logger {
	return &ZapLog{
		logger: l.logger.With(zap.Any(key, value)),
	}
}

func (l *ZapLog) WithError(err error) Logger {
	return &ZapLog{
		logger: l.logger.With(zap.Any("error", err.Error())),
	}
}

func (l *ZapLog) Error(msg string) {
	l.logger.Error(msg)
}

func (l *ZapLog) Info(msg string) {
	l.logger.Info(msg)
}

func (l *ZapLog) Debug(msg string) {
	l.logger.Debug(msg)
}
